//container elements

//cluster info

//apply elements
let applyFile = document.getElementById("applyFile");
let applySubmit = document.getElementById("applySubmit");
let applyForm = document.getElementById("applyForm");

//delete elements
let deleteSelect = document.getElementById("deleteSelect");
let deleteName = document.getElementById("deleteName");
let deleteSubmit = document.getElementById("deleteSubmit");
let deleteForm = document.getElementById("deleteForm");

//get elements
let getSelect = document.getElementById("getSelect");
let getName = document.getElementById("getName");
let getSubmit = document.getElementById("getSubmit");
let getForm = document.getElementById("getForm");

//describe elements
let describeSelect = document.getElementById("describeSelect");
let describeName = document.getElementById("describeName");
let describeSubmit = document.getElementById("describeSubmit");
let descibeForm = document.getElementById("describeForm");

//scale elements
let replicaCount = document.getElementById("replicaCount");
let replicaSubmit = document.getElementById("replicaSubmit");
let scaleForm = document.getElementById("scaleForm");

//log elements
let logName = document.getElementById("logName");
let logSubmit = document.getElementById("logSubmit");
let logForm = document.getElementById("logForm");

//form links
let showFormLinks = document.getElementsByClassName("showForm");

//controller address
let server = "http://localhost:4040";

document.addEventListener("DOMContentLoaded", () => {
    console.log("All content loaded");
});

function submitApply(filepath) {
    let reqBody = {
        filepath: filepath,
        object: "",
        name: "",
        replicas: "",
    };
    fetch(`${server}/apply`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    });
}

function submitGet(object, name) {
    let reqBody = {
        filepath: "",
        object: object,
        name: name,
        replicas: "",
    };
    console.log(`${server}/get`);
    return fetch(`${server}/get`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    })
    .then(resp => resp.text())
    .then(json => {
        return json;
    });
}

function deleteObj(object, name) {
    let reqBody = {
        filepath: "",
        object: object,
        name: name,
        replicas: "",
    };
    fetch(`${server}/delete`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    });
}

function submitReplica(count) {
    let reqBody = {
        filepath: "",
        object: "",
        name: "",
        replicas: count,
    };
    fetch(`${server}/scale`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    });
}

function submitDescribe(object, name) {
    let reqBody = {
        filepath: "",
        object: object,
        name: name,
        replicas: "",
    }
    return fetch(`${server}/describe`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    })
    .then(resp => resp.text())
    .then(json => {
        return json;
    });
}

function getLogs(name) {
    let reqBody = {
        filepath: "",
        object: "",
        name: name,
        replicas: "",
    }
    return fetch(`${server}/logs`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    })
    .then(resp => resp.text())
    .then(json => {
       return json;
    });
}

function showForm(id) {
    el = document.getElementById(id);
    el.style.display = "flex";
}

function hideForm(id) {
    el = document.getElementById(id);
    el.style.display = "none";
}

function hideAllForms() {
    applyForm.style.display = "none";
    getForm.style.display = "none";
    describeForm.style.display = "none";
    scaleForm.style.display = "none";
    deleteForm.style.display = "none";
    logForm.style.display = "none";
}

for (let link of showFormLinks) {
    link.addEventListener("click", e => {
        hideAllForms();
        showForm(e.target.value);
    });
}

applySubmit.addEventListener("click", () => {
    submitApply(applyFile.value);
    hideForm("applyForm");
});

getSubmit.addEventListener("click", () => {
    submitGet(getSelect.value, getName.value).then(got => {
        console.log(got);
    });
    hideForm("getForm");
});

deleteSubmit.addEventListener("click", () => {
    deleteObj(deleteSelect.value, deleteName.value);
    hideForm("deleteForm");
});

replicaSubmit.addEventListener("click", () => {
    submitReplica(replicaCount.value);
    hideForm("scaleForm");
});

describeSubmit.addEventListener("click", () => {
    submitDescribe(describeSelect.value, describeName.value).then(desc => {
        console.log(desc);
    });
    hideForm("describeForm");
});

logSubmit.addEventListener("click", () => {
    getLogs(logName.value).then(log => {
        console.log(log);
    });
    hideForm("logForm");
});