(() => {
  for (const quiz of document.querySelectorAll("[data-quiz]")) {
    const feedback = quiz.querySelector("[data-feedback]");
    const answers = quiz.querySelectorAll("[data-answer]");

    for (const answer of answers) {
      answer.addEventListener("click", () => {
        for (const option of answers) {
          option.classList.remove("is-correct", "is-incorrect");
          option.setAttribute("aria-pressed", "false");
        }

        const correct = answer.dataset.correct === "true";
        answer.classList.add(correct ? "is-correct" : "is-incorrect");
        answer.setAttribute("aria-pressed", "true");

        if (feedback) {
          feedback.textContent = answer.dataset.feedback || "";
          feedback.className = `feedback ${correct ? "good" : "bad"}`;
        }
      });
    }
  }
})();
