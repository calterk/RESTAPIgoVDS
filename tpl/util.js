/*function Parent() {
    this.name = 'parent';

    this.nameStr = () => {
        return this.name;
    };
}*/

function Util() {
    let token = '';
    const ls = localStorage.getItem('token');
    if (ls !== null) token = ls;
    this.setToken = tk => {
        token = tk;
    }

    this.post = (url, data, callback) => {
        fetch(url, {
            method: 'POST',
            body: JSON.stringify(data),
            headers: {
                'Content-type': 'application/json',
                'Authorization': token,
            },
        })
            .then(data => data.json())
            .then(callback);
    };

    this.get = (url, callback) => {
        fetch(url, {
            method: 'GET',
            headers: {
                'Content-type': 'application/json',
                'Authorization': token,
            },
        })
            .then(data => data.json())
            .then(callback);
    };

    this.id = (el) => document.getElementById(el);
    this.$ = (el) => document.querySelector(el);
    this.$$ = (el) => document.querySelectorAll(el);
    this.modals = {};

    this.modal = (id, action) => {
        if (!this.modals[id]) {
            this.modals[id] = new bootstrap.Modal(document.getElementById(id));
        }

        this.modals[id][action]();
    };
    /*
    this.parse = (content, params) => {
        let param = Object.assign({}, params);

        return content.replace(/{{(\w+)}}/g, (str) => {
            str = str.substring(2, str.length - 2);
            if (param[str] === undefined) return '';
            return param[str];
        });
    };

    this.tpl = `<table class="table">
    <tr>
        <th>uid</th>
        <th>login</th>
        <th>pass</th>
        <th>name</th>
        <th>access</th>
    </tr>
    {{users}}
</table>`;

    this.tr = `<tr>
        <td>{{uid}}</td>
        <td>{{login}}</td>
        <td>{{pass}}</td>
        <td>{{name}}</td>
        <td>{{access}}</td>

        </td>
    </tr>`;*/
}

const util = new Util();
