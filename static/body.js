
///// THEMING ////
// Set cookie
function setTheme(theme) {
    if (theme === "auto") {
        document.documentElement.removeAttribute("data-theme");
    } else {
        document.documentElement.setAttribute("data-theme", theme);
    }
    document.cookie = `theme=${theme};path=/;max-age=31536000`; // 1 year
}

// Read cookie
function getTheme() {
    const match = document.cookie.match(/theme=(\w+)/);
    return match ? match[1] : "light";
}

// Apply on page load
const saved = getTheme();
setTheme(saved);
document.getElementById("theme-switcher").value = saved;

// Listen for changes
document.getElementById("theme-switcher").addEventListener("change", (e) => {
    setTheme(e.target.value);
});

//// LANGUAGES ////
document.getElementById("language-switcher").addEventListener("change", (e => {
    document.cookie = `lang=${e.target.value};path=/;max-age=31536000`;
    location.reload()
}))

document.querySelectorAll('input[name="source"]').forEach(el => {
    el.addEventListener("change", () => {
        const checked = [...document.querySelectorAll('input[name="source"]:checked')]
            .map(cb => cb.value)
            .join(",");
        document.cookie = `sources=${checked};path=/;max-age=31536000`;
        htmx.trigger(document.getElementById("articles"), "refreshArticles");
    });
});

document.querySelectorAll('input[name="source"]').forEach(el => {
    el.addEventListener("change", () => {
    });
});
