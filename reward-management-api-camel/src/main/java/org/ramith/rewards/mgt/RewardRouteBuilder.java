package org.ramith.rewards.mgt;

import org.apache.camel.builder.RouteBuilder;
import org.apache.camel.model.rest.RestBindingMode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

@Component
public class RewardRouteBuilder extends RouteBuilder {

    private Logger  log = LoggerFactory.getLogger(RewardRouteBuilder.class);
    public void configure() throws Exception {
        restConfiguration().component("servlet").bindingMode(RestBindingMode.json);

        rest("/select-reward").post().type(RewardSelection.class).outType(String.class)
                .to("direct:handleRewardSelection");

        from("direct:handleRewardSelection")
                .process(exchange -> {
                    RewardSelection selection = exchange.getIn().getBody(RewardSelection.class);
                    exchange.setProperty("selection", selection);
                })
                .toD("netty-http:${header.LOYALTY_API_URL}/user/${body.userId}")
                .process(exchange -> {
                    User user = exchange.getIn().getBody(User.class);
                    RewardSelection selection = exchange.getProperty("selection", RewardSelection.class);
                    Reward reward = transform(user, selection);
                    exchange.getIn().setBody(reward);
                })
                .toD("netty-http:${header.VENDOR_MANAGEMENT_API_URL}/rewards")
                .setBody(constant("reward selection received successfully"));

    }

    private Reward transform(User user, RewardSelection rewardSelection) {
        Reward reward = new Reward(rewardSelection.selectedRewardDealId(), user.userId(), user.firstName(), user.lastName(), user.email());
        return reward;
    }

}



record User(String userId, String firstName, String lastName, String email) {}
record RewardSelection(String userId, String selectedRewardDealId, boolean acceptedTnC) {}
record Reward(String rewardId, String userId, String firstName, String lastName, String email) {}

