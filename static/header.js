const m = document.cookie.match(/theme=(\w+)/);
if (m && m[1] !== "auto") {
    document.documentElement.setAttribute("data-theme", m[1]);
}

document.querySelectorAll('input[name="source"]').forEach(el => {
    el.addEventListener("change", () => {
        const checked = [...document.querySelectorAll('input[name="source"]:checked')]
            .map(cb => cb.value)
            .join(",");
        document.cookie = `sources=${checked};path=/;max-age=31536000`;
        htmx.trigger("#articles", "load");
    });
});