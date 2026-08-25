"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getAllGroups, createGroup } from "@/lib/api/groups";
import { useGroupChat } from "@/contexts/GroupChatContext";
import GroupCard from "@/components/groups/GroupCard";
import NailButton from "@/components/shared/NailButton";
import "@/css/groups.css";
import { useToast } from "@/contexts/ToastContext";

export default function GroupsPage() {
  const router = useRouter();
  const { unread } = useGroupChat();
  const [groups, setGroups] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [creating, setCreating] = useState(false);
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [submitting, setSubmitting] = useState(false);
  //add tooast to grouppp creation 
  const  toooasst=useToast()

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const response = await getAllGroups();
        if (!cancelled) setGroups(response?.data || []);
      } catch (err) {
        
        if (
          err.status === 401 ||
          err.status === 403 ||
          err.status === 404 ||
          err.status >= 500
        ) {
          router.push(
            `/error?message=${encodeURIComponent(err.statusText)}`
          );
        }
        if (!cancelled) setError(err?.message || "Could not load groups.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, []);

  async function handleCreate(event) {
    event.preventDefault();
    setError("");
    setSubmitting(true);
    try {
      const response = await createGroup({
        title: title.trim(),
        description: description.trim(),
        is_public: 1,
      });
      const group = response?.data;
      if (group?.id) {
        setTitle("");
        setDescription("");
        setCreating(false);
        router.push(`/group/${group.id}`);
      } else {
        setError("Unexpected response from the server.");
      }
    } catch (err) {
       
        if (
          err.status === 401 ||
          err.status === 403 ||
          err.status === 404 ||
          err.status >= 500
        ) {
          router.push(
            `/error?message=${err.statusText}`
          );
        }
      // setError(err?.message || "Could not create group.");
          toooasst.error(err?.message ||"Could not create group.")
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section className="postsContainer">
      <div className="feedHeader">
        <div>
          <h1 className="feedTitle">Groups</h1>
        </div>
        <NailButton onClick={() => setCreating(!creating)}>
          {creating ? "Close" : "Create group"}
        </NailButton>
      </div>

      {creating && (
        <form className="groupForm" onSubmit={handleCreate}>
          {error && <div className="postsError">{error}</div>}

          <div className="field">
            <label htmlFor="group-title">Title</label>
            <input
              id="group-title"
              value={title}
              onChange={(event) => setTitle(event.target.value)}
              placeholder="Group title"
              required
            />
          </div>

          <div className="field">
            <label htmlFor="group-description">Description</label>
            <textarea
              id="group-description"
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              required
            />
          </div>

          <div className="postFormActions">
            <button className="button postFormButton" type="submit" disabled={submitting}>
              {submitting ? "Creating..." : "Create group"}
            </button>
          </div>
        </form>
      )}

      {loading && <div className="postsPlaceholder">Loading groups...</div>}
      {error && !creating && <div className="postsError">{error}</div>}
      {!loading && !error && groups.length === 0 && (
        <div className="postsPlaceholder">No groups yet. Create the first one.</div>
      )}

      <div className="groupList">
        {groups.map((group) => (
          <GroupCard key={group.id} group={group} unreadCount={unread[group.id] || 0} />
        ))}
      </div>
    </section>
  );
}
