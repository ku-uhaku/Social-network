"use client";

import { useState, useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { createPost } from "@/lib/api/posts";
import { getFollowers } from "@/lib/api/user";
import ImageUploadButton from "@/components/shared/ImageUploadButton";
import UsersSelect from "@/components/shared/UsersSelect";
import { useAuth } from "@/contexts/AuthContext";
import "@/css/createPost.css";
import { useToast } from "@/contexts/ToastContext";

export default function CreatePostPage() {
  const router = useRouter();
  const { user: currentUser } = useAuth();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [privacy, setPrivacy] = useState("public");
  const [image, setImage] = useState(null);
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [followers, setFollowers] = useState([]);
  const [selectedViewers, setSelectedViewers] = useState([]);
  const [loadingFollowers, setLoadingFollowers] = useState(false);
  const toooasst=useToast()
  console.log("that is the toastt",toooasst)

  const params = useSearchParams();
  const groupId = params.get("group_id") ? Number(params.get("group_id")) : null;
  // Fetch current user followers when privacy is private
  useEffect(() => {
    if (privacy !== "private" || !currentUser?.id) return;

    let cancelled = false;
    getFollowers(currentUser.id)
      .then((response) => {
        if (!cancelled) setFollowers(response.data || []);
      })
      .catch(() =>
        toooasst.error("Failed to load followers list")
      )
      .finally(() => {
        if (!cancelled) setLoadingFollowers(false);
      });
    return () => {
      cancelled = true;
    };
  }, [privacy, currentUser?.id]);

  async function handleSubmit(event) {
    event.preventDefault();
    // setError("");
    setSubmitting(true);

    if (title.trim() === "" || content.trim() === "") {
      toooasst.error("Title and content cannot be empty or contain only spaces");
      setSubmitting(false);
      return;
    }

    try {
      // Validate private posts
      if (privacy === "private" && selectedViewers.length === 0) {
        // setError("Private posts must specify at least one viewer");
            toooasst.error(err?.message ||"Private posts must specify at least one viewer");
        setSubmitting(false);
        return;
      }

      const formData = new FormData();
      formData.append("title", title.trim());
      formData.append("content", content.trim());
      console.log("groupId:::,!!!::", groupId);
      if (groupId) {
        formData.append("group_id", groupId);
        formData.append("privacy", "group");
      } else {
        formData.append("privacy", privacy);
      }
      if (privacy === "private") {
        selectedViewers.forEach((id) => formData.append("visible_to", id));
      }
      if (image) formData.append("image", image);

      const response = await createPost(formData);
      const post = response?.data;
      if (post?.id) {
        router.push(groupId ? `/group/${groupId}` : `/`);
      } else {
        // setError("Unexpected response from the server.");
            toooasst.error("Unexpected response from the server.");
      }
    } catch (err) {
      // setError(err?.message || "Could not create post.");
            toooasst.error(err?.message || "Could not create post.");
    } finally {
      setSubmitting(false);
    }
  }


  return (
    <section className="postsContainer">
      <div className="postFormWrapper">
        <h1 className="feedTitle">Create a new post</h1>
        <form className="postForm" onSubmit={handleSubmit}>
          {error && <div className="postsError">{error}</div>}

          <div className="field">
            <label htmlFor="title">Title</label>
            <input
              id="title"
              value={title}
              onChange={(event) => setTitle(event.target.value)}
              placeholder="Post title"
              required
              minLength={1}
              maxLength={40}
            />
          </div>

          <div className="field">
            <label htmlFor="content">Content</label>
            <textarea
              id="content"
              value={content}
              onChange={(event) => setContent(event.target.value)}
              placeholder="Write your post here..."
              required
              maxLength={200}
            />
          </div>

          <div className="row">
            
              {groupId ? (<></>) : (
                <div className="field">
                  <label htmlFor="privacy">Privacy</label>
                  <select
                    id="privacy"
                    value={privacy}
                    onChange={(event) => {
                      const next = event.target.value;
                      setPrivacy(next);
                      setLoadingFollowers(next === "private");
                      if (next !== "private") {
                        setSelectedViewers([]);
                        setFollowers([]);
                      }
                    }}
                  >
                    <option value="public">Public</option>
                    <option value="almost private">Almost Private</option>
                    <option value="private">Private</option>
                  </select>
                </div>
              )}
            
            <div className="field">
              <ImageUploadButton
                label="Image (optional)"
                value={image}
                onChange={setImage}
              />
            </div>
          </div>


          {/* Follower selection for private posts */}
          {privacy === "private" && (
            <div className="field">
              <label>Select Viewers</label>
              <UsersSelect
                users={followers}
                loading={loadingFollowers}
                selected={selectedViewers}
                onToggle={(id, checked) => {
                  if (checked) {
                    setSelectedViewers([...selectedViewers, id]);
                  } else {
                    setSelectedViewers(
                      selectedViewers.filter((viewerID) => viewerID !== id)
                    );
                  }
                }}
              />
              {selectedViewers.length === 0 && (
                <div className="errorText">
                  At least 1 person
                </div>
              )}
            </div>
          )}

          <div className="postFormActions">
            <button
              className="button postFormButton"
              type="submit"
              disabled={submitting || (privacy === "private" && selectedViewers.length === 0)}
            >
              {submitting ? "Creating..." : "Publish post"}
            </button>
          </div>
        </form>
      </div>
    </section>
  );
}
