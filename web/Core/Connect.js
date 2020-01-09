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
    }

    /**
     * Send line to ws server
     *
     * @param {Array} points
     * @returns {Boolean}
     */
    sendLine(points) {
        if (this._socket.readyState == this._socket.OPEN) {
            this._socket.send(points);
            return true;
        }
        return false;
    }
}