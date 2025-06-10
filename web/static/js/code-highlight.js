import { createHighlighter } from "https://esm.sh/shiki@3.0.0";

async function highlightCode() {
  const highlighter = await createHighlighter({
    themes: ["github-light"],
    langs: ["go", "javascript", "typescript", "yaml", "bash", "dockerfile"],
  });

  document.querySelectorAll("pre code").forEach((block) => {
    const lang = block.className.replace("language-", "") || "text";
    if (highlighter.getLoadedLanguages().includes(lang)) {
      console.log("executing highlighter", lang);
      const html = highlighter.codeToHtml(block.textContent, {
        lang: lang,
        theme: "github-light",
      });
      block.parentElement.outerHTML = html;
    }
  });
}

highlightCode();
