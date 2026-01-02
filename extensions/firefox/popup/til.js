// Selectors
let message = document.querySelector("#message");
let loginPage = document.querySelector("#login-page");
let addContentPage = document.querySelector("#add-content");

// Check current storage to define which screen to display
browser.storage.local.get("feedKey")
    .then(res => {
        if (!res.feedKey) {
            throw new Error("No feed key found");
        }
        addContentPage.classList.remove("hidden");
        loginPage.classList.add("hidden");
    })
    .catch(() => {
        addContentPage.classList.add("hidden");
        loginPage.classList.remove("hidden");
    })

browser.tabs.query({active: true, currentWindow: true})
    .then(tabs => {
        document.querySelector("#title").value = tabs[0].title;
        document.querySelector("#link").value = tabs[0].url;
    });


// Listeners
document.querySelector("#login-form").addEventListener("submit", login);
document.querySelector("#add-content-form").addEventListener("submit", addContent);
document.querySelector("#logout-button").addEventListener("click", logout);

// Set listener on message div so we can remove it automatically after 5s
let currentTimeout = null;
let callback = () => {
    if (currentTimeout) {
        clearTimeout(currentTimeout);
    }

    currentTimeout = setTimeout(() => {
        message.innerHTML = "";
        message.classList.forEach(c => message.classList.remove(c));
    }, 5000)
};

let observer = new MutationObserver(callback);
observer.observe(message, {attributes: true, childList: true});


// Utility functions

function login() {
    event.preventDefault();
    setMessage("Checking credentials...", "info")
    loginPage.classList.add("hidden");

    // Try to get feed. If valid, set information. Otherwise, display error.
    let serverUrl = document.querySelector("#server_url").value;
    let feedKey = document.querySelector("#feed_key").value;

    if (serverUrl.endsWith("/")) {
        serverUrl = serverUrl.slice(0, -1);
    }

    fetch(`${serverUrl}/api/rss?key=${feedKey}`)
        .then(res => {
            if (!res.ok) {
                setMessage("Invalid credentials.", "error");
                loginPage.classList.remove("hidden");
                return
            }
            setMessage("Connected!", "success");
            browser.storage.local.set({
                server: serverUrl,
                feedKey: feedKey
            })
            addContentPage.classList.remove("hidden");
        })
        .catch(() => {
            setMessage("Invalid credentials.", "error");
            loginPage.classList.remove("hidden");
        })
}

function addContent() {
    event.preventDefault();
    setMessage("Adding content...", "info");
    document.querySelector("#add-content-form").classList.add("hidden");

    let tags = document.querySelector("#tags").value.split(",");
    tags.push(`lang:${document.querySelector("#lang").value}`)

    let form = {
        title: document.querySelector("#title").value,
        link: document.querySelector("#link").value,
        tags: tags,
        content: document.querySelector("#content").value
    }

    browser.storage.local.get().then(storage => {
        fetch(`${storage.server}/api/posts?key=${storage.feedKey}`, {
            method: "POST",
            body: JSON.stringify(form),
        })
            .then(res => {
                if (!res.ok) {
                    res.text().then(t => {
                        setMessage(`Error adding content: ${t}`, "error");
                    })
                } else {
                    setMessage("Content added!", "success");
                    document.querySelector("#title").value = "";
                    document.querySelector("#link").value = "";
                    document.querySelector("#tags").value = "";
                    document.querySelector("#content").value = "";
                }
            })
            .catch(err => {
                setMessage(`Error adding content: ${err}`, "error");
            })
            .finally(() => {
                document.querySelector("#add-content-form").classList.remove("hidden");
            })
    })

}

function logout() {
    browser.storage.local.clear();
    setMessage("Logged out.", "info");
    loginPage.classList.remove("hidden");
    addContentPage.classList.add("hidden");
}

function setMessage(displayMessage, cssClass) {
    message.innerHTML = `<p>${displayMessage}</p>`;
    message.classList.forEach(c => message.classList.remove(c));
    message.classList.add(cssClass);
}