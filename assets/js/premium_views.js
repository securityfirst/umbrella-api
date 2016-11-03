function PremiumsCommitLedgerView() {
    this.base = View; this.base(); var self = this;

    self.init = function() {
        self.$sel(".refresh-btn").click(function() {
            self.goto("loading");
            self.postNotification("premiums_ledger_commit:"+self.context.uid);
        });
    }

    self.onLoading = function() {
        $.ajax({
            url: "/data/get_user/" + self.context.uid + "/premiums_ledger_commits",
        }).done(function(res) {
            if (res.length > 0) {
                self.context.commits = res;
                self.goto("main")
            } else {
                self.goto("no-records")
            }
        })
    };

    self.onMain = function() {
        var $table = self.$actionSel().find(".table");

        $table.html("");

        self.context.commits.forEach(function(e) {
            // Insert row
            embed("PremiumsCommitLedgerRowView", $table, {
                commit: e,
                uid: self.context.uid
            });
        });
    };
}

function PremiumsCommitLedgerRowView() {
    this.base = View; this.base(); var self = this;

    self.init = function() {
        self.$sel(".toggle").on("click", function() {
            self.context.updates = self.context.updates || {};

            self.context.updates['hidden'] = !$(this).is(':checked')
            self.goto("update", $(this).val())
        })
    }

    self.onLoading = function() {
        $.ajax({
            url: "/data/premiums_ledger_commit/" + self.context.commit.id + "/get_object",
        }).done(function(res) {
            self.context.commitObject = res;
            self.goto("main");
        });
    };

    self.onUpdate = function() {
        $.ajax({
            method: "put",
            url: "/data/premiums_ledger_commit/" + self.context.commit.id,
            contentType: 'application/json',
            data: JSON.stringify(self.context.updates),
        }).done(function(res) {
            self.context.updates = {};
            self.context.commit = res;
            self.postNotification("premiums_ledger_commit:"+self.context.uid);
            self.goto("main")
        });
    }

    self.onMain = function() {
        var $content = self.$actionSel().find('.table-row-content');
        var $rowEdge = self.$actionSel().find(".table-row-edge")
        var $description = $content.find(".description")
        var $toggle = $content.find(".toggle")
        var premiumTransaction = self.context.premiumTransaction;
        var commit = self.context.commit;

        // Make sure toggle is in the correct state
        $toggle.prop("checked", !commit.hidden)

        // Dim disabled
        if (commit.hidden == true) {
            self.$actionSel().css('opacity', 0.3);
        } else {
            self.$actionSel().css('opacity', 1.0);
        }

        // Change color on the row based on what it is
        $rowEdge.removeClass('row-edge-expired')
        $rowEdge.removeClass('row-edge-active')

        if (commit.expired == true) {
            $rowEdge.addClass('row-edge-expired')
        } else {
            $rowEdge.addClass('row-edge-active')
        }

        if (commit.type === 'premium_transaction') {
            $description.html("");
            embed("PremiumsCommitLedgerRowPremiumTransactionPartial", $description, {
                premium_transaction: self.context.commitObject,
                commit: self.context.commit,
            });
        }

        var $icon = self.$actionSel().find('.table-row-edge .icon');
        $icon.addClass('fa-credit-card');
    };
}

function PremiumsStatsView() {
    this.base = View; this.base(); var self = this;

    this.init = function() {
      self.onNotification("premiums_ledger_commit:"+self.context.uid, function() {
        self.goto("loading")
      });
    }

    this.onLoading = function() {
        $.ajax({
            url: "/data/user/" + self.context.uid + "/payment_stats",
        }).done(function(res) {
            self.context.paymentStats = res;
            self.goto("main");
        });
    }

    this.onMain = function() {
        self.$actionSel().find(".usd").html(self.context.paymentStats.estimated_spent_usd)
    }
}

function PremiumsCommitLedgerRowPremiumTransactionPartial() {
    this.base = View; this.base(); var self = this;

    self.onLoading = function() {
        // Show apple store style transaction
        if (self.context.premium_transaction.is_apple_store === true) {
            self.goto("apple_store")
        }
    }

    self.onAppleStore = function() {
        var commit = self.context.commit;
        var expiresOn = commit.expires_on;
        var createdAt = commit.created_at;
        var name = self.context.premium_transaction.price_name;
        console.log(self.context);
        var estimatedUSD = self.context.premium_transaction.estimated_usd;
        self.$actionSel().find(".estimated-usd .usd").html(estimatedUSD);

        // Created at
        self.$actionSel().find(".created-time").html(moment.unix(createdAt).format("MM/DD/YYYY"));

        // Reset expired
        self.$actionSel().find(".expire-time").removeClass('expired')
        self.$actionSel().find(".expire-time").removeClass('not-expired')

        if (commit.does_expire === true) {
            self.$actionSel().find(".expire-time time").removeClass('hidden')
            self.$actionSel().find(".expire-time-na").addClass('hidden')
            self.$actionSel().find(".expire-time time").timeago("update", (new Date(expiresOn * 1000)))

            if (commit.expired === true) {
                self.$actionSel().find(".expire-time").addClass('expired')
            } else {
                self.$actionSel().find(".expire-time").addClass('not-expired')
            }
        } else {
            self.$actionSel().find(".expire-time").removeClass('expired')
            self.$actionSel().find(".expire-time").addClass('not-expired')

            self.$actionSel().find(".expire-time time").addClass('hidden')
            self.$actionSel().find(".expire-time-na").removeClass('hidden')
        }

        self.$actionSel().find(".name").html(name);
    }
}

function PremiumFeaturesView() {
    this.base = View; this.base(); var self = this;

    self.init = function() {
      self.onNotification("premiums_ledger_commit:"+self.context.uid, function() {
        self.goto("loading")
      });

      console.log(self)
    }

    self.onLoading = function() {
        $.ajax({
            url: "/data/get_user/" + self.context.uid + "/premium_features",
        }).done(function(res) {
            self.context.premiumFeatures = res;
            self.goto("main");
        });
    };

    self.onMain = function() {
        self.$actionSel().html("")
        function insertFeature(feature) {
            embed('PremiumFeatureCardView', self.$actionSel(), {
                feature: feature
            });
        }


        $.each(self.context.premiumFeatures, function(i, e) {
            insertFeature(e);
        });
    }
}

function PremiumFeatureCardView() {
    this.base = View; this.base(); var self = this; 
    
    self.onMain = function() {
        var feature = self.context.feature;

        self.$actionSel().find(".card-title h1").html(feature.name);

        if (feature.active == false) {
            self.$actionSel().find('.card-big-value .checked').removeClass('checked')
        }
    };
}
