//buttons
let editButton = document.getElementById("editButton");
let addButton = document.getElementById("addButton");
let getButton = document.getElementById("getButton");

//editor area
let editorArea = document.getElementById("editorArea");

//user login
let userLogin = document.getElementById("userLogin");
let userCreate = document.getElementById("userCreate"); 
let loginFrame = document.getElementById("loginFrame");
let usernameInput = document.getElementById("username");
let passwordInput = document.getElementById("password");
let createUserInput = document.getElementById("createUsername");
let createPasswordInput = document.getElementById("createPassword");
let loginErr = document.getElementById("loginErr");
let createUserErr = document.getElementById("createUserErr");

//webpage header
let userMenu = document.getElementById("userMenu");
let loginNav = document.getElementById("loginNav");
let userObj = {};

let loggedIn = false;

//change this variable to point to the address of the server
let authServer = "http://localhost:5555";

//once the page has loaded, check if the user is logged in already via session/cookies
document.addEventListener("DOMContentLoaded", () => {
    fetch(authServer + "/login").then(resp => {
        if (resp.status == 200) {
            loggedIn = true;
            resp.json().then(json => {
                userObj = json;
            });
        }
        else {
            showLoginScreen();
        }
    });
});

//doesnt do anything yet
editButton.addEventListener("click", () => {
    console.log("Sorry, the edit button hasn't been implemented yet");
    // fetch(authServer + "/editJournal").then(resp => {
    //     if (resp != 200) {
    //         showLoginScreen();
    //     }
    //     return resp.json();
    // }).then(json => {
    //     console.log(json);
    // });
});

//gets the user's entire journal and puts the latest into the editor
getButton.addEventListener("click", () => {
    fetch(authServer + "/getJournal", {
        method: "POST",
        headers: {"Content-Type": "text/plain"},
        body: JSON.stringify({
            Username: userObj.username,
            Journal: [],
        })
    }).then(resp => {
        if (resp.status != 200) {
            showLoginScreen();
        }
        return resp.json();
    }).then(json => {
        let output = "";
        json.journal.forEach(entry => {
            output += `${entry.date}\n${entry.entry}\n\n\n`
        });
        editorArea.value = output;
    });
});

//adds a journal entry to the db for today's date
addButton.addEventListener("click", () => {
    fetch(authServer + "/addJEntry", {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify({
            Username: userObj.username,
            Journal: [
                {
                    date: "",
                    entry: editorArea.value
                }
            ],
        })
    }).then(resp => {
        if (resp.status != 200) {
            showLoginScreen();
        }
    });
});

userLogin.addEventListener("click", e => {
    e.preventDefault();
    let name = usernameInput.value;
    let pass = passwordInput.value;
    if (validateUser(name) && validatePass(pass)) {
        login(name, pass);
        hideLoginScreen();
        clearLogins();
    } else {
        loginErr.innerText = "Invalid character in username or password, try again";
    }
});

userCreate.addEventListener("click", e => {
    e.preventDefault();
    let name = createUserInput.value;
    let pass = createPasswordInput.value;
    if (validateUser(name) && validatePass(pass)) {
        createUser(name, pass).then(created => {
            if (created) {
                login(name, pass);
            }
        });
        hideLoginScreen();
        clearLogins();
    } else {
        createUserErr.innerText = "Invalid character in username or password, try again";
    }
});

//shows login screen when user clicks "log in" button
loginNav.addEventListener("click", () => {
    showLoginScreen();
});

//clears login input fields when user clicks or tabs into login window
loginFrame.addEventListener("focus", () => {
    loginErr.innerText = "";
    createUserErr.innerText = "";
});

showLoginScreen = () => {
    loginFrame.style.zIndex = 2;
    loginFrame.style.display = "block";
}

hideLoginScreen = () => {
    loginFrame.style.zIndex = -1;
    loginFrame.style.display = "none";
}

//checks user credentials from db and returns a boolean indicating
//if the user was successfully logged in or not
let login = (user, pass) => {
    return fetch(authServer + "/login", {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify({
            username: user,
            password: pass
        })
    }).then(resp => {
        if (resp.status == 200) {
            resp.json().then(json => {
                userObj = json;
                updateUserMenu();
            });
        }
        loggedIn = resp.status == 200;
        return resp.status == 200;
    });
}

//creates a new user in the db and returns a boolean
//indicating whether user creation was successful
createUser = (user, pass) => {
    return fetch(authServer + "/createUser", {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify({
            username: user,
            password: pass
        })
    }).then(resp => {
        return resp.status == 200;
    });
}

//updates the user menu in navigation with the current contents of userObj
let updateUserMenu = () => {
    if (loggedIn) {
        userMenu.innerHTML = `<button id="userDropdown" class="navButton">&#9660; ${userObj.username}</button>`;
    }
    else {
        userMenu.innerHTML = `<button id="loginNav" class="navButton">Log In</button>`
    }
}

let clearLogins = () => {
    usernameInput.value = "";
    passwordInput.value = "";
    createUserInput.value = "";
    createPasswordInput.value = "";
}

let validateUser = user => user.search(/[ !@#$%^&()\[\]{}`~:;<>,.\/\\+*"?']/) == -1;

let validatePass = pass => pass.search(/[ ()\[\]{}~:;<>,.\/\\+"']/) == -1;