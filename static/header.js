const m = document.cookie.match(/theme=(\w+)/);
if (m && m[1] !== "auto") {
    document.documentElement.setAttribute("data-theme", m[1]);
}