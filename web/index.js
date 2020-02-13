//buttons
let editButton = document.getElementById("editButton");
let addButton = document.getElementById("addButton");

//editor area
let editorArea = document.getElementById("editorArea");

//user login
let userLogin = document.getElementById("userLogin");
let userCreate = document.getElementById("userCreate"); 
let loginFrame = document.getElementById("loginFrame");
let usernameInput = document.getElementById("username");
let passwordInput = document.getElementById("password");
let createUserInput = document.getElementById("createUsername");
let createEmailInput = document.getElementById("createPassword");
let createPasswordInput = document.getElementById("createEmail");

//webpage header
let userMenu = document.getElementById("userMenu");
let loginNav = document.getElementById("loginNav");
let userObj = {};

let loggedIn = false;

//once the page has loaded, check if the user is logged in already
document.addEventListener("DOMContentLoaded", () => {
    fetch("/login").then(resp => {
        if (resp.status == 200) {
            loggedIn = true;
            user = resp.json().then(json => userObj = json);
        }
        else {
            showLoginScreen();
        }
    });
});

editButton.addEventListener("click", () => {
    fetch("/editJournal").then(resp => {
        if (resp != 200) {
            showLoginScreen();
        }
        return resp.json();
    }).then(json => {
        console.log(json);
    });
});

addButton.addEventListener("click", () => {
    fetch("/addJournal", {
        method: "post",
        body: {
            date: "",
            text: editorArea.nodeValue,
        }
    }).then(resp => {
        if (resp != 200) {
            showLoginScreen();
        }
        return resp.json();
    }).then(json => {
        console.log(json);
    });
});

userLogin.addEventListener("click", e => {
    e.preventDefault();
    let name = usernameInput.value;
    let pass = passwordInput.value;
    login(name, pass).then(res => {
        loggedIn = res;
    });
    hideLoginScreen();
    clearLogins();
});

userCreate.addEventListener("click", e => {
    e.preventDefault();
    let email = createEmailInput.value;
    let name = createUserInput.value;
    let pass = createPasswordInput.value;
    createUser(name, pass, email).then(res => {
        loggedIn = res;
    });
    hideLoginScreen();
    clearLogins();
});

loginNav.addEventListener("click", () => {
    showLoginScreen();
});

showLoginScreen = () => {
    loginFrame.style.zIndex = 2;
    loginFrame.style.display = "block";
}

hideLoginScreen = () => {
    loginFrame.style.zIndex = -1;
    loginFrame.style.display = "none";
}

let login = (user, pass) => {
    return fetch("/login", {
        method: "post",
        body: {
            user: user,
            pass: pass,
        }
    }).then(resp => {
        if (resp.status == 200) {
            loggedIn = true;
            resp.json().then(json => userObj = json);
            updateUserMenu();
        }
        return resp.status == 200;
    });
}

createUser = (user, pass, email) => {
    return fetch("/createUser", {
        method: "post",
        body: {
            user: user,
            pass: pass,
            email: email,
        }
    }).then(resp => {
        return resp.status == 200;
    });
}

let updateUserMenu = () => {
    if (loggedIn) {
        userMenu.innerHTML = `<button id="userDropdown">${user.name}</button>`;
    }
    else {
        userMenu.innerHTML = `<button id="loginNav">Log In</button>`
    }
}

let clearLogins = () => {
    usernameInput.value = "";
    passwordInput.value = "";
    createUserInput.value = "";
    createEmailInput.value = "";
    createPasswordInput.value = "";
}
