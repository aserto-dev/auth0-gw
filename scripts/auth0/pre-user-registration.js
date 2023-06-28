/**
* Handler that will be called during the execution of a PreUserRegistration flow.
*
* @param {Event} event - Details about the context and user that is attempting to register.
* @param {PreUserRegistrationAPI} api - Interface whose methods can be used to change the behavior of the signup.
*/
exports.onExecutePreUserRegistration = async (event, api) => {
    const { httpTransport, emitterFor, CloudEvent } = require("cloudevents");

    // Create an emitter to send events to a receiver
    const emit = emitterFor(httpTransport("http://auth0-gw.eastus.cloudapp.azure.com:8383/events"));

    const type = 'pre-user-registration';
    const source = 'auth0.com';
    const data = event.user;

    // Create a new CloudEvent
    const ce = new CloudEvent({ type, source, data });

    // Send it to the endpoint - encoded as HTTP binary by default
    emit(ce);
};
