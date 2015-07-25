var FileUpload = React.createClass({
    // when a file is passed to the input field, retrieve the contents as a
    // base64-encoded data URI and send the object to the this.props.onData func
    onChange: function(e) {
        console.log('onChange:', e);
        var reader = new FileReader();
        var file = e.target.files[0];

        reader.onload = function (upload) {
            console.log('upload.onload', upload);
            this.props.onData(upload.target.result);
        }.bind(this);

        //reader.readAsBinaryString(file);
        reader.readAsDataURL(file);
    },
    onSubmit: function(event) {
        console.log('onSubmit: ', event);
        event.preventDefault();
    },
    render: function() {
        return (
            <form onSubmit={this.onSubmit} encType="multipart/form-data">
                <input type="file" onChange={this.onChange} />
            </form>
        )
    }
});
var Uploader = React.createClass({
    getInitialState: function() {
        return {
            data: null,
            loading: false,
            failed: false,
            translated: null,
        };
    },
    shipIt: function(b64str) {
        var win = function (data, status, xhr) {
            console.log("win", data);
            this.setState({loading: false, translated: data});
        }.bind(this);
        var fail = function(xhr, status, err) {
            console.log("error", err);
            this.setState({loading: false, failed: true});
        }.bind(this);
        $.ajax({
            url: '/post/new',
            method: 'POST',
            data: '{"images": ["'+b64str+'"]}',
            contentType: 'application/json',
            dataType: 'text', // we expect plaintext back
            success: win,
            error: fail,
        })
    },
    onData: function (data) {
        console.log('onData:', data);
        var split = data.match(/,(.*)$/)[1];
        this.shipIt(split);
        this.setState({
            loading: true,
            b64String: split,
            data: data,
        });
    },
    render: function() {
        return (
            <section>
                <FileUpload onData={this.onData} />
                <section className="preview">
                    { this.state.translated ? <pre>{this.state.translated}</pre> : null }
                    { this.state.data       ? <img src={this.state.data} />      : <p>Preview will be here</p> }
                    { this.state.failed     ? <p>'Error translating image.'</p>  : null }
                </section>
            </section>
        )
    }
});
