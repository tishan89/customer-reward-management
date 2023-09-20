import ballerina/http;
import ballerina/log;
import ballerina/oauth2;


public type RewardSelection record {
    string userId;
    string selectedRewardDealId;
    boolean acceptedTnC;
};

public type User record {
    string userId;
    string firstName;
    string lastName;
    string email;
};


public type Reward record {
    string rewardId;
    string userId;
    string firstName;
    string lastName;
    string email;
};

configurable string clientId = ?;
configurable string clientSecret = ?;
configurable string tokenUrl = ?;
configurable string loyaltyApiUrl = ?;
configurable string vendorManagementApiUrl = ?;

oauth2:ClientOAuth2Provider provider = new ({
    tokenUrl: tokenUrl,
    clientId: clientId,
    clientSecret: clientSecret
});

# A service representing a network-accessible API
# bound to port `9090`.
service / on new http:Listener(9090) {

    resource function post pick(RewardSelection selection) returns error|string {
        log:printInfo("reward selected: ", selection = selection);

        http:Client clientEp = check new (loyaltyApiUrl, {
            auth: {
                tokenUrl: tokenUrl,
                clientId: clientId,
                clientSecret: clientSecret
            }
        });
        User|http:Error user = clientEp->/user/[selection.userId];
        if (user is http:Error) {
            log:printError("error retrieving user: ", 'error = user);
            return user;
        }

        log:printInfo("user retrieved: ", user = user);
        Reward reward = transform(user, selection);

        http:Client vendorManagementClientEp = check new (vendorManagementApiUrl, {
            auth: {
                tokenUrl: tokenUrl,
                clientId: clientId,
                clientSecret: clientSecret
            }
        });

        http:Response|http:Error response = vendorManagementClientEp->post("/rewards", reward);

        if response is http:Error {
            log:printError("error while sending reward selection to vender ", 'error = response);
            return response;
        }
        
        log:printInfo("reward selection sent to vendor ", statusCode = response.statusCode);
        return "success";

    }

}

function transform(User user, RewardSelection rewardSelection) returns Reward => {

    firstName: user.firstName,
    lastName: user.lastName,
    userId: rewardSelection.userId,
    rewardId: rewardSelection.selectedRewardDealId,
    email: user.email
};
