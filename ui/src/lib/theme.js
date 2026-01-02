// Init on browser load
if(!localStorage.getItem('theme')) {
    localStorage.setItem('theme', "dark")
}

if(localStorage.getItem('theme') === 'white') {
    addWhiteTheme()
}

export function changeTheme() {
    if(localStorage.getItem('theme') === 'dark') {
        localStorage.setItem('theme', "white")
        addWhiteTheme()
    } else {
        localStorage.setItem('theme', "dark")
        removeWhiteTheme()
    }
}

function addWhiteTheme() {
    let head  = document.getElementsByTagName('head')[0];
    let link  = document.createElement('link');
    link.rel  = 'stylesheet';
    link.type = 'text/css';
    link.id = 'white-theme';
    link.href = "/style-white.css";
    head.appendChild(link);
}

function removeWhiteTheme() {
    document.getElementById('white-theme').remove();
}

