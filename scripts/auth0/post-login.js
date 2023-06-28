/**
* Handler that will be called during the execution of a PostLogin flow.
*
* @param {Event} event - Details about the user and the context in which they are logging in.
* @param {PostLoginAPI} api - Interface whose methods can be used to change the behavior of the login.
*/
exports.onExecutePostLogin = async (event, api) => {
    const { httpTransport, emitterFor, CloudEvent } = require("cloudevents");

    // Create an emitter to send events to a receiver
    const emit = emitterFor(httpTransport("http://auth0-gw.eastus.cloudapp.azure.com:8383/events"));

    const type = 'post-login';
    const source = 'auth0.com';
    const data = event.user;
    
    // Create a new CloudEvent
    const ce = new CloudEvent({ type, source, data });

    // Send it to the endpoint - encoded as HTTP binary by default
    emit(ce);
};


/**
* Handler that will be invoked when this action is resuming after an external redirect. If your
* onExecutePostLogin function does not perform a redirect, this function can be safely ignored.
*
* @param {Event} event - Details about the user and the context in which they are logging in.
* @param {PostLoginAPI} api - Interface whose methods can be used to change the behavior of the login.
*/
// exports.onContinuePostLogin = async (event, api) => {
// };
