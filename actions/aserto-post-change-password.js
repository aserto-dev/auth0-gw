/**
* Handler that will be called during the execution of a PostChangePassword flow.
*
* @param {Event} event - Details about the user and the context in which the change password is happening.
* @param {PostChangePasswordAPI} api - Methods and utilities to help change the behavior after a user changes their password.
*/
exports.onExecutePostChangePassword = async (event, api) => {
    const { httpTransport, emitterFor, CloudEvent } = require("cloudevents");

    // Create an emitter to send events to a receiver
    const emit = emitterFor(httpTransport("http://<auth0-gw-hostname>:8383/events"));

    const type = 'post-change-password';
    const source = 'auth0.com';
    const data = event.user;

    // Create a new CloudEvent
    const ce = new CloudEvent({ type, source, data });

    // Send it to the endpoint - encoded as HTTP binary by default
    emit(ce);
};
