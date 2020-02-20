//container elements
let deploymentSection = document.getElementById("deploymentSection");
let deploymentCardContainer = document.getElementById("deploymentCardContainer");
let nodeSection = document.getElementById("nodeSection");
let nodeCardContainer = document.getElementById("nodeCardContainer");
let podSection = document.getElementById("podSection");
let podCardContainer = document.getElementById("podCardContainer");
let serviceSection = document.getElementById("serviceSection");
let serviceCardContainer = document.getElementById("serviceCardContainer");
let podMenuScale = document.getElementById("podMenuScale");

//cluster info
let clusterInfo = document.getElementById("clusterInfo");

//apply elements
let applyFile = document.getElementById("applyFile");
let applySubmit = document.getElementById("applySubmit");
let applyForm = document.getElementById("applyForm");

//run elements
let runForm = document.getElementById("runForm");
let runImage = document.getElementById("runImage");
let runSubmit = document.getElementById("runSubmit");

//log elements
let logName = document.getElementById("logName");
let logForm = document.getElementById("logForm");

//form links
let showFormLinks = document.getElementsByClassName("showForm");

//controller address
let server = "http://localhost:4040";

//kuberenetes objects
let deployment = [];
let nodes = [];
let pods = [];
let services = [];

document.addEventListener("DOMContentLoaded", () => {
    console.log("All content loaded");

    //get cluster info and add it into the page
    fetch(`${server}/cluster`).then(resp => resp.text())
    .then(text => {
        //remove all the junk from the outputted string
        text = text.replace(/\[0;32m/gi, "");
        text = text.replace(/\[0m/gi, "");
        text = text.replace(/\[0;33m/gi, "");
        text = text.split("\n")

        //add the text to H3 tags and put it into the HTML
        html = `<h3>${text[0]}</h3><h3>${text[1]}</h3>`;
        clusterInfo.innerHTML += html;
    })

    //Refresh the cluster object info at the given interval
    setInterval(() => {
        submitGet("deployment", "").then(json => {
            deployment = json.items;
            html = ``;
            podMenuScale.innerHTML = "";
            deployment.forEach(dep => {
                html += buildDeploymentCard(dep);
                addDeploymentToScale(dep.metadata.name, dep.status.replicas);
            });
            deploymentCardContainer.innerHTML = "";
            deploymentCardContainer.innerHTML += html;
        });
        submitGet("node", "").then(json => {
            nodes = json.items;
            html = ``;
            nodes.forEach(node => {
                html += buildNodeCard(node);
            });
            nodeCardContainer.innerHTML = "";
            nodeCardContainer.innerHTML += html;
        });
        submitGet("pod", "").then(json => {
            pods = json.items;
            html = ``;
            pods.forEach(pod => {
                html += buildPodCard(pod);
            });
            podCardContainer.innerHTML = "";
            podCardContainer.innerHTML += html;
        });
        submitGet("svc", "").then(json => {
            services = json.items;
            html = ``;
            services.forEach(service => {
                html += buildServiceCard(service);
            });
            serviceCardContainer.innerHTML = "";
            serviceCardContainer.innerHTML += html;
        });
    }, 3000);
});

function buildDeploymentCard(deployment) {
    //create container list
    containers = ``;
    deployment.spec.template.spec.containers.forEach(container => {
        containers += `<p>Container Name: ${container.name}</p>
                       <p>Container Image: ${container.image}</p>`
    });

    html = `<div class="objectCard">
                <h2>${deployment.metadata.name}</h2>
                <p>Namespace: ${deployment.metadata.namespace}</p>
                <p>Replicas: ${deployment.status.replicas}</p>
                <p>Available Replicas: ${deployment.status.availableReplicas}</p>
                ${containers}
                <ul>
                    <button onclick="deleteObj('deployment', '${deployment.metadata.name}')"><img src="./trash.svg" alt="Delete Deployment"></button>
                </ul>
            </div>`;
    return html;
}

function buildServiceCard(service) {
    //create port list
    ports = ``;
    service.spec.ports.forEach(port => {
        ports += `<p>Port: ${port.protocol} ${port.port}</p>
                  <p>Target Port: ${port.protocol} ${port.targetPort}</p>
                  <p>Node Port: ${port.protocol} ${port.nodePort}</p>`;
    });

    html = `<div class="objectCard">
                <h2>${service.metadata.name}</h2>
                <p>Service Type: ${service.spec.type}</p>
                <p>Namespace: ${service.metadata.namespace}</p>
                <p>Internal IP: ${service.spec.clusterIP}</p>
                ${ports}
                <ul>
                    <button onclick="deleteObj('svc', '${service.metadata.name}')"><img src="./trash.svg" alt="Delete Deployment"></button>
                </ul>
            </div>`;
    return html;
}

function buildNodeCard(node) {
    //create condition message list
    condition = ``;
    node.status.conditions.forEach(con => {
        condition += `<p>${con.message}</p>`;
    });

    //create addresses list
    addresses = ``;
    node.status.addresses.forEach(addr => {
        addresses += `<p>${addr.type}: ${addr.address}</p>`;
    });

    html = `<div class="objectCard">
                <h2>${node.metadata.name}</h2>
                ${addresses}
                ${condition}
            </div>`;

    return html;
}

function buildPodCard(pod) {
    //create container list
    containers = ``;
    pod.spec.containers.forEach(container => {
        containers += `<p>Container Name: ${container.name}</p>
                       <p>Container Image: ${container.image}</p>`
    });

    html = `<div class="objectCard">
                <h2>${pod.metadata.name}</h2>
                <p>Namespace: ${pod.metadata.namespace}</p>
                <p>On Node: ${pod.spec.nodeName}</p>
                ${containers}
                <ul>
                    <button onclick="showLogs('${pod.metadata.name}')"><img src="./logs.svg" alt="View Logs"></button>
                    <button onclick="deleteObj('pod', '${pod.metadata.name}')"><img src="./trash.svg" alt="Delete Deployment"></button>
                </ul>
            </div>`;
    return html;
}

function submitApply(yaml) {
    let reqBody = {
        deployment: yaml,
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
        deployment: "",
        object: object,
        name: name,
        replicas: "",
    };
    return fetch(`${server}/get`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    })
    .then(resp => resp.json())
    .then(json => {
        return json;
    });
}

function deleteObj(object, name) {
    let reqBody = {
        deployment: "",
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

function submitReplica(count, deployment) {
    let reqBody = {
        deployment: `deployment/${deployment}`,
        object: "",
        name: "",
        replicas: count + "",
    };
    fetch(`${server}/scale`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
    });
}

function submitRun(image) {
    let reqBody = {
        deployment: "",
        object: "",
        name: image,
        replicas: "",
    };
    fetch(`${server}/run`, {
        method: "POST",
        headers: { "Content-Type": "text/plain" },
        body: JSON.stringify(reqBody),
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
    .then(text => {
       return text;
    });
}

function addDeploymentToScale(name, replicas) {
    let html = `<h2>${name}</h2><button name="${name}" onclick="submitReplica('${replicas+1}', '${name}')">&Delta;</button><h3>${replicas}</h3><button name="${name}" onclick="submitReplica('${replicas-1}', '${name}')">&nabla;</button>`;
    podMenuScale.innerHTML += html;
}

function showLogs(name) {
    getLogs(name).then(logs => {
        showForm("logForm");
        logDisplay.value = logs;
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
    runForm.style.display = "none";
    logForm.style.display = "none";
}

for (let link of showFormLinks) {
    link.addEventListener("click", e => {
        hideAllForms();
        showForm(e.target.value);
    });
}

applySubmit.addEventListener("click", () => {
    let yaml = document.getElementById("deploymentScript");
    submitApply(yaml.value);
    hideForm("applyForm");
});

runSubmit.addEventListener("click", () => {
    let image = runImage.value;
    submitRun(image);
    hideForm("runForm");
});