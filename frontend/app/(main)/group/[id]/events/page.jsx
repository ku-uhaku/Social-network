"use client";
import { useEffect, useRef, useState } from "react";
import { useParams, notFound } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useToast } from "@/contexts/ToastContext";
import { getGroup, getGroupEvents, createGroupEvent, cancelGroupEvent, setEventResponse } from "@/lib/api/groups";
import NailButton from "@/components/shared/NailButton";
import "@/css/groups.css";
import { useRouter } from "next/navigation";

export default function GroupEventsPage() {
  const { id } = useParams();
  const { user: currentUser } = useAuth();
  const toast = useToast();
  const groupId = Number(id);
  if (isNaN(groupId)) notFound();
    // const { unreadCount } = useNotifications();


  const [group, setGroup] = useState(null);
  const [membership, setMembership] = useState("none");
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [formOpen, setFormOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [dateTime, setDateTime] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const pendingRef = useRef(new Set());

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const g = await getGroup(groupId);
        if (cancelled) return;
        setGroup(g?.data?.group || null);
        setMembership(g?.data?.membership || "none");
        // backend enforces member-only access
        const ev = await getGroupEvents(groupId);
        if (!cancelled) setEvents(ev?.data || []);
      } catch (err) {
      
        if (err?.message === "Group not found") notFound();
        if (err?.message === "you do not have permission to access or post to this context") {
          setMembership("none");
        } else if (!cancelled) {
          setError(err?.message || "Could not load events.");
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [groupId]);

  function reloadEvents() {
    getGroupEvents(groupId).then((ev) => setEvents(ev?.data || []));
  }

  async function handleCreate(e) {
    e.preventDefault();
    if (!dateTime) { toast.error("Choose a date and time for the event."); return; }
    if (title.trim() === "" || description.trim() === "") {
      toast.error("Title and description cannot be empty or contain only spaces");
      return;
    }
    setSubmitting(true);
    try {
      await createGroupEvent({ group_id: groupId, title, description, event_time: new Date(dateTime).toISOString() });
      setTitle(""); setDescription(""); setDateTime(""); setFormOpen(false);
      reloadEvents();
    } catch (err) {
        
      toast.error(err?.message || "Could not create event.", { action: { label: "Go Home", href: "/" } });
    }
    finally { setSubmitting(false); }
  }

  function applyLocalResponse(ev, status) {
    if (ev.my_status === status) return ev;
    let going = ev.going_count;
    let notGoing = ev.not_going_count;
    if (ev.my_status === "going") going -= 1;
    if (ev.my_status === "not_going") notGoing -= 1;
    if (status === "going") going += 1;
    if (status === "not_going") notGoing += 1;
    return { ...ev, my_status: status, going_count: going, not_going_count: notGoing };
  }

  async function handleChoice(eventId, status) {
    if (pendingRef.current.has(eventId)) return;
    pendingRef.current.add(eventId);
    const prev = events.find((ev) => ev.id === eventId);
    if (!prev) {
      pendingRef.current.delete(eventId);
      return;
    }

    setEvents((list) =>
      list.map((ev) => (ev.id === eventId ? applyLocalResponse(ev, status) : ev))
    );

    try {
      const res = await setEventResponse(eventId, status);
      const updated = res?.data;
      if (updated?.id === eventId) {
        setEvents((list) => list.map((ev) => (ev.id === eventId ? updated : ev)));
      }
    } catch (err) {
        
      setEvents((list) => list.map((ev) => (ev.id === eventId ? prev : ev)));
      toast.error(err?.message || "Could not update your response.");
    } finally {
      pendingRef.current.delete(eventId);
    }
  }

  async function handleCancel(eventId) {
    try { await cancelGroupEvent(eventId); reloadEvents(); }
    catch (err) {
        toast.error(err?.message || "Could not cancel event."); }
  }

  if (loading) return <div className="postsPlaceholder">Loading events...</div>;
  if (error) return <div className="postsError">{error}</div>;
  if (!group) notFound();

  const member = membership === "accepted";

  return (
    <section className="postsContainer">
      <div className="groupHeader">
        <h1 className="feedTitle">{group.title} - Events</h1>
        {member && (
          <button className="eventMakeButton" type="button" onClick={() => setFormOpen(!formOpen)}>
            {formOpen ? "Close form" : "Create event"}
          </button>
        )}
      </div>

      {!member ? (
        <div className="postsPlaceholder">Join this group to see its events.</div>
      ) : (
        <>
          {formOpen && (
            <form className="groupForm" onSubmit={handleCreate}>
              <label>Title <input value={title} onChange={(e) => setTitle(e.target.value)} required maxLength={20} /></label>
              <label>Description <textarea value={description} onChange={(e) => setDescription(e.target.value)} required maxLength={100} /></label>
              <label>Date / Time <input type="datetime-local" value={dateTime} onChange={(e) => setDateTime(e.target.value)} required /></label>
              <NailButton type="submit" disabled={submitting}>{submitting ? "Creating..." : "Create event"}</NailButton>
            </form>
          )}

          {events.length === 0 ? (
            <div className="postsPlaceholder">No events yet.</div>
          ) : (
            <div className="eventList">
              {events.map((ev) => {
                const isOver = ev.status !== "upcoming";
                const isCreator = ev.creator_id === currentUser?.id;
                return (
                  <div key={ev.id} className={`eventCard ${isOver ? "eventExpired" : ""}`}>
                    <div className="eventCardHeader">
                      <h2 className="eventCardTitle">{ev.title}</h2>
                      <span className="eventStatus">{ev.status}</span>
                    </div>
                    {ev.creator_username && (
                      <div className="eventCardCreator">Created by {ev.creator_username}</div>
                    )}
                    <p className="eventCardDescription">{ev.description}</p>
                    <div className="eventCardTime">{new Date(ev.event_time).toLocaleString()}</div>
                    {!isOver && (
                      <div className="eventChoices">
                        <button type="button" className={`eventChoice ${ev.my_status === "going" ? "active" : ""}`} onClick={() => handleChoice(ev.id, "going")}>Going ({ev.going_count})</button>
                        <button type="button" className={`eventChoice ${ev.my_status === "not_going" ? "active" : ""}`} onClick={() => handleChoice(ev.id, "not_going")}>Not going ({ev.not_going_count})</button>
                        {isCreator && <button type="button" className="eventCancelButton" onClick={() => handleCancel(ev.id)}>Cancel event</button>}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </>
      )}
    </section>
  );
}