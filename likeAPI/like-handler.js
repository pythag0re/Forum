document.addEventListener('DOMContentLoaded', () => {
  const getCookie = (name) => {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    return parts.length === 2 ? parts.pop().split(';').shift() : null;
  };

  const userId = getCookie('user_id');
  if (!userId) {
    console.warn("Utilisateur non connecté.");
    return;
  }

  document.querySelectorAll('.post').forEach(postElement => {
    const postId = postElement.dataset.postId;
    const likeBtn = postElement.querySelector('.like-btn');
    const likeCount = postElement.querySelector('.like-count');

    if (!likeBtn || !likeCount || !postId) return;

    fetch(`http://localhost:3001/like/has-liked?user_id=${userId}&post_id=${postId}`)
      .then(res => res.json())
      .then(data => {
        if (data.liked) {
          likeBtn.innerHTML = '<i class="fa-solid fa-thumbs-up"></i> Unlike';
          likeBtn.dataset.liked = "true";
        } else {
          likeBtn.dataset.liked = "false";
        }
      });

    likeBtn.addEventListener('click', () => {
      const liked = likeBtn.dataset.liked === "true";
      const method = liked ? 'DELETE' : 'POST';

      fetch('http://localhost:3001/like', {
        method: method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: parseInt(userId),
          post_id: parseInt(postId)
        })
      })
      .then(() => {
        const count = parseInt(likeCount.textContent);
        likeCount.textContent = liked ? count - 1 : count + 1;
        likeBtn.innerHTML = liked
          ? '<i class="fa-regular fa-thumbs-up"></i> Like'
          : '<i class="fa-solid fa-thumbs-up"></i> Unlike';
        likeBtn.dataset.liked = liked ? "false" : "true";
      });
    });
  });
});
