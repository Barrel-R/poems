// TODO: set axios

function getPoems() {
    const url = "http://localhost:8000/api/v1/poems"

    axios.get(url)
        .then(response => {
            return response.data
        })
        .catch(errors => {
            console.log(errors)
        })
}

function getPoem(id) {
    const url = "http://localhost:8000/api/v1/poems"

    axios.get(url + "/" + id)
        .then(response => {
            console.log(response.data)
        })
        .catch(errors => {
            console.log(errors)
        })
}
