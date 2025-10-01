function save_post(id, title, markdown, html) {
    console.log('id:',      id)
    console.log('title:',   title)
    console.log('md:',      markdown)
    console.log('html:',    html)

    fetch('/api/post/save', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            id: id,
            title: title,
            markdown: markdown,
            html: html
        })
    })
    .then(res => res.text())
    .then(data => console.log('resp: ', data))
    .catch(err => console.log('err:', err))
}
