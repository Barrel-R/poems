async function getPoems() {
    const url = "http://localhost:8080/api/v1/poems"
    const errorBag = {
        message: "",
        status: null,
    }

    try {
        const response = await fetch(url)

        if (!response.ok) {
            errorBag.message = "The server returned an error: " + response.status
            errorBag.status = response.status

        }
        const json = await response.json()

        return json

    } catch (error) {
        errorBag.status = error.status
        errorBag.message = error.message
        console.error(error)
        console.error(errorBag)
    }
}

async function getPoem(id) {
    const url = "http://localhost:8080/api/v1/poems/" + id
    const errorBag = {
        message: "",
        status: null,
    }

    try {
        const response = await fetch(url)

        if (!response.ok) {
            errorBag.message = "The server returned an error: " + response.status
            errorBag.status = response.status
        }

        const json = await response.json()

        return json

    } catch (error) {
        errorBag.status = error.status
        errorBag.message = error.message
        console.error(error)
        console.error(errorBag)

        return error
    }
}

function setPoems(data) {
    const list = document.getElementById("poemList")

    for (let i = 0; i < data.length; i++) {
        const item = document.createElement("div")
        const title = document.createElement("h3")
        const textDisplay = document.createElement("div")

        list.appendChild(item)
        item.appendChild(title)
        item.appendChild(textDisplay)

        item.classList = "poem-item"
        title.classList = "poem-title crimson"
        textDisplay.classList = "poem crimson-i"

        stanza = data[i]?.content.split("\n\n")
        title.innerText = data[i].title
        textDisplay.innerText = stanza[0] + "\n..."

        console.log(data[i])
    }
}

function displayPoem(data) {
    const poemSection = document.getElementById("poem")

    const poem = document.createElement("div")
    const title = document.createElement("h3")
    const textDisplay = document.createElement("div")

    poemSection.appendChild(poem)
    poem.appendChild(title)
    poem.appendChild(textDisplay)

    poem.classList = "single-poem"
    title.classList = "poem-title crimson"
    textDisplay.classList = "poem crimson-i"

    title.innerText = data.title
    textDisplay.innerText = data.content

    console.log(data)
}
