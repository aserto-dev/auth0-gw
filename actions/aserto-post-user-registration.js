/**
* Handler that will be called during the execution of a PostUserRegistration flow.
*
* @param {Event} event - Details about the context and user that has registered.
* @param {PostUserRegistrationAPI} api - Methods and utilities to help change the behavior after a signup.
*/
exports.onExecutePostUserRegistration = async (event, api) => {
    const { httpTransport, emitterFor, CloudEvent } = require("cloudevents");

    // Create an emitter to send events to a receiver
    const emit = emitterFor(httpTransport("http://<auth0-gw-hostname>:8383/events"));

    const type = 'post-user-registration';
    const source = 'auth0.com';
    const data = event.user;

    // Create a new CloudEvent
    const ce = new CloudEvent({ type, source, data });

    // Send it to the endpoint - encoded as HTTP binary by default
    emit(ce);
};
