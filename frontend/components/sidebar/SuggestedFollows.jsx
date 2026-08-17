"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Avatar from "@/components/shared/Avatar";
import { getSuggestedUsers, followUser } from "@/lib/api/user";

export default function SuggestedFollows() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [busyId, setBusyId] = useState(null);

  useEffect(() => {
    //tooo check if the component is mounted
    let cancelled = false;
    async function load() {
      try {
        const response = await getSuggestedUsers(5);
        const data = response?.data || response || [];
        if (!cancelled) setUsers(Array.isArray(data) ? data : []);
      } catch (err) {
        //theee error defiandeed right nooow 
        if (!cancelled) setError(err?.message || "Could not load suggestions.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
  }, []);

  async function handleFollow(userId) {
    //take the id 
    console.log("take the id",userId)
    if (busyId) return;
    setBusyId(userId);
    try {
      await followUser(userId);
      // Drrrop the follllowed user soo that means the list stay frech
      // console.log("thossse is user suggested")
      setUsers((prev) => prev.filter((u) => u.id !== userId));
    } catch {
      ///machi dabaaaa
    } finally {
      setBusyId(null);
    }
  }

  if (loading || error || users.length === 0) return null;
  console.log("the users ready to display right now ",users)

  return (
    <aside className="suggestedFollows">
      <h2 className="suggestedFollowsTitle">Suggested for you</h2>
      <div className="suggestedFollowsCard">
        {users.map((u) => (
          <div key={u.id} className="suggestedFollowItem">
            <Avatar avatar={u.avatar} username={u.username} size={40} />
            <Link href={`/profile/${u.username}`} className="suggestedFollowInfo">
              <span className="suggestedFollowName">
                {u.first_name} {u.last_name}
              </span>
              <span className="suggestedFollowUsername">@{u.username}</span>
            </Link>
            <button
              type="button"
              className="suggestedFollowButton"
              onClick={() => handleFollow(u.id)}
              disabled={busyId === u.id}
            >
              {busyId === u.id ? "..." : "Follow"}
            </button>
          </div>
        ))}
      </div>
    </aside>
  );
}
