import Line from "./Line.js";

export default class Connect {
    /**
     * Connect constructor
     *
     * @param {String} host
     * @param {Function} callback
     */
    constructor(host, callback) {
        this._socket = new WebSocket(`ws://${host}/ws`);
        this._socket.onmessage = (event) => {
            let data = JSON.parse(event.data);
            callback(data);
        };
        this._stack = [];
    }

    /**
     * Send line to ws server
     *
     * @param {Line} line
     * @returns {Boolean}
     */
    sendLine(line) {
        this._stack.push(line);
        return this;
    }

    /**
     * Send stack to server
     */
    send() {
        if (this._socket.readyState == this._socket.OPEN) {
            this._socket.send(points);
            return true;
        }
        return false;
    }
}