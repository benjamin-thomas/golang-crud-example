/**
 * This file provided by Facebook is for non-commercial testing and evaluation
 * purposes only. Facebook reserves all rights not expressly granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * FACEBOOK BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 * WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

var Comment = React.createClass({
  render: function() {
    var rawMarkup = marked(this.props.children.toString(), {sanitize: true});
    return (
      <div className="comment">
        <h2 className="commentAuthor">
          {this.props.author}
        </h2>
        <span dangerouslySetInnerHTML={{__html: rawMarkup}} />
      </div>
    );
  }
});

var CommentBox = React.createClass({
  loadCommentsFromServer: function() {
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  handleCommentSubmit: function(comment) {
    var comments = this.state.data;
    var newComments = comments.concat([comment]);
    this.setState({data: newComments});
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      type: 'POST',
      data: comment,
      success: function(data) {
        this.setState({data: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount: function() {
    this.loadCommentsFromServer();
    setInterval(this.loadCommentsFromServer, this.props.pollInterval);
  },
  render: function() {
    return (
      <div className="commentBox">
        <h1>Comments</h1>
        <CommentList data={this.state.data} />
        <CommentForm onCommentSubmit={this.handleCommentSubmit} />
      </div>
    );
  }
});

var CommentList = React.createClass({
  render: function() {
    var commentNodes = this.props.data.map(function(comment, index) {
      return (
        // `key` is a React-specific concept and is not mandatory for the
        // purpose of this tutorial. if you're curious, see more here:
        // http://facebook.github.io/react/docs/multiple-components.html#dynamic-children
        <Comment author={comment.author} key={index}>
          {comment.text}
        </Comment>
      );
    });
    return (
      <div className="commentList">
        {commentNodes}
      </div>
    );
  }
});

var CommentForm = React.createClass({
  handleSubmit: function(e) {
    e.preventDefault();
    var author = React.findDOMNode(this.refs.author).value.trim();
    var text = React.findDOMNode(this.refs.text).value.trim();
    if (!text || !author) {
      return;
    }
    this.props.onCommentSubmit({author: author, text: text});
    React.findDOMNode(this.refs.author).value = '';
    React.findDOMNode(this.refs.text).value = '';
  },
  render: function() {
    return (
      <form className="commentForm" onSubmit={this.handleSubmit}>
        <input type="text" placeholder="Your name" ref="author" />
        <input type="text" placeholder="Say something..." ref="text" />
        <input type="submit" value="Post" />
      </form>
    );
  }
});

var DataTable = React.createClass({
    loadData: function(q) {
        var params = '';
        if (q) {
            params += '?name=' + 'per=10&page=1&q=%25'+q+'%25&cols=name&condOp=OR&matchOp=ILIKE';
        }
        var fullURL = this.props.url+params;
        console.log("fullURL:", fullURL);
        $.ajax({
            url: fullURL,
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({data: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(this.props.url, status, err.toString());
            }.bind(this)
        });
    },

    getInitialState: function() {
        return {
            data: [],
            cols: {},
            params: '?per=10&page=1',
        };
    },

    componentDidMount: function() {
        this.loadData();
    },

    onUserInput: function(search, name) { // searchValue, searchColumn
        console.log("search again:", search, "name:", name);
        this.loadData(search);

        var cols = this.state.cols;
        if (search === "") {
            delete cols[name];
        } else {
            cols[name] = search;
        }
        this.setState({
            cols: cols,
        })
        var keys = Object.keys(cols);
        params = '?per=10&page=1';
        for (var i = 0; i < keys.length; i++) {
            var k = keys[i];
            params += "&" + k + "=" + cols[k];
        }
        this.setState({
            params: params,
        })
    },

    render: function() {
        return (
            <div className="responsive-table">
              <div className="scrollable-area">
                <table className="table table-bordered table-striped">
                    <TableHeaders onUserInput={this.onUserInput} />
                    <TableRows data={this.state.data} />
                </table>

                <p>
                  Selected: {JSON.stringify(this.state.cols)}
                </p>

                <p>
                URL: {this.state.params}
                </p>

              </div>
            </div>
        );
    }
});

var TableHeaders = React.createClass({
    onUserInput: function(search, name) {
        console.log("search in headers: ", search, "name:", name);
        this.props.onUserInput(search, name);
    },
    render: function() {
        var attrs = [
            { label: 'ID',       name: 'id'},
            { label: 'Name',     name: 'name'},
            { label: 'Line1',    name: 'line1'},
            { label: 'Line2',    name: 'line2'},
            { label: 'Line3',    name: 'line3'},
            { label: 'City',     name: 'city'},
            { label: 'Zip Code', name: 'zip code'},
            { label: 'Country',  name: 'country'}
        ],
        headers = [],
        searchers = [];

        for (var i = 0; i < attrs.length; i++) {
            attr = attrs[i];
            headers.push(<TableHeader value={attr.label} />);
            searchers.push(<TableSearcher name={attr.name} onUserInput={this.onUserInput} />);
        }
        return (
            <thead>
                <tr>
                    {headers}
                </tr>
                <tr>
                    {searchers}
                </tr>
            </thead>
        );
    }
});

var TableHeader = React.createClass({
    getInitialState: function() {
        return {
            count: 0,
        };
    },
    onClick: function() {
        this.setState({count: this.state.count+1});
    },
    render: function() {
        var dir = ['ASC', 'DESC'][this.state.count%2]
        return (
            <th className="no-mouse-select" onClick={this.onClick}>{this.props.value} [{dir}] ({this.state.count})</th>
        );
    }
});

var TableSearcher = React.createClass({
    onKeyUp: function(e) {
        var search = e.target.value;
        this.props.onUserInput(search, this.props.name)
    },
    render: function() {
        return (
            <th>
                <input type="text" name={this.props.name} onKeyUp={this.onKeyUp} />
            </th>
        )
    }
});

var TableRows = React.createClass({
    getInitialState: function() {
        return {
            data: [],
        }
    },
    render: function() {
        var rows = this.props.data.map(function(addr) {
            return (
                <TableRow data={addr} />
            );
        });
        return (
            <tbody>
                {rows}
            </tbody>
        );
    }
});

var TableRow = React.createClass({
    render: function() {
        var addr = this.props.data;
        var tds = Object.keys(addr).map(function(k, v) {
            return (
                <td>{addr[k]}</td>
            )
        })
        return (
            <tr>
                {tds}
            </tr>
        )
    }
});

// React.render(
//   <CommentBox url="comments.json" pollInterval={2000} />,
//   document.getElementById('content')
// );
React.render(
    <DataTable url="/api/addresses" />,
    document.getElementById('content')
);
