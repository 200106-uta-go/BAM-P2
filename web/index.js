//buttons
let editButton = document.getElementById("editButton");
let addButton = document.getElementById("addButton");

//editor area
let editorArea = document.getElementById("editorArea");

editButton.addEventListener("click", () => {
    fetch("/edit").then(resp => resp.json()).then(json => {
        console.log(json);
    });
});

addButton.addEventListener("click", () => {
    fetch("/GoJournal", {
        method: "post",
        body: {
            date: ""
        }
    }).then(resp => resp.json()).then(json => {
        console.log(json);
    });
});
