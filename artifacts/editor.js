function save_post(id, title, markdown, html, alternative, tags) {
    console.log('id:',          id)
    console.log('title:',       title)
    console.log('md:',          markdown)
    console.log('html:',        html)
    console.log('alternative:', alternative)
    console.log('tags:',        tags)

    fetch_with_json('/api/post/save',
        {
            id: id,
            title: title,
            markdown: markdown,
            html: html,
            alternative: alternative,
            tags: tags
        }
    )
}

function fetch_with_json(route, obj, callback = null) {
    fetch(route, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(obj)
    })
    .then(res => res.text())
    .then(data => {
        console.log('resp: ', data)
        if (null != callback) {
            try {
                callback(JSON.parse(data))
            }
            catch (e) {
                callback(data)
            }
        }
    })
    .catch(err => console.log('err:', err))
}
