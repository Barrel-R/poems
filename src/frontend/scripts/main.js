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

        item.id = data[i].id
        stanza = data[i]?.content.split("\n\n")
        title.innerText = data[i].title
        textDisplay.innerText = stanza[0] + "\n..."

        item.addEventListener("click", (event) => {
            window.location.href = "/poem.html?poem=" + data[i].id.toString()
        })

        console.log(data[i])
    }
}

function displayPoem(response) {
    const poemSection = document.getElementById("poem")
    const pagination = document.getElementById("pagination")

    const poem = document.createElement("div")
    const title = document.createElement("h3")
    const textDisplay = document.createElement("div")
    const currentPoemId = document.createElement("a")
    const paginationNextLink = document.createElement("a")
    const paginationPrevLink = document.createElement("a")

    if (response.pagination.previousPoemId != 0) {
        paginationPrevLink.href = "/poem.html?poem=" + response.pagination.previousPoemId
        paginationPrevLink.classList += "left-link"
    } else {
        paginationPrevLink.setAttribute("disabled", "")
    }

    poemSection.appendChild(poem)
    poem.appendChild(title)
    poem.appendChild(textDisplay)
    pagination.appendChild(paginationPrevLink)
    pagination.appendChild(currentPoemId)
    pagination.appendChild(paginationNextLink)

    if (response.pagination.nextPoemId != 0) {
        paginationNextLink.href = "/poem.html?poem=" + response.pagination.nextPoemId
        paginationNextLink.classList += "right-link"
    } else {
        paginationNextLink.setAttribute("disabled", "")
    }

    poem.classList = "single-poem"
    title.classList = "poem-title crimson"
    textDisplay.classList = "poem crimson-i"

    title.innerText = response.data.title
    textDisplay.innerText = response.data.content
    currentPoemId.innerText = response.data.id
    paginationPrevLink.innerHTML = `&laquo`
    paginationNextLink.innerHTML = `&raquo`
    currentPoemId.classList += "active"
    currentPoemId.href = "#"
}
