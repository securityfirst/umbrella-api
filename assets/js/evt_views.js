function EvtEndpointsListView() {
    this.base = View; this.base(); var self = this; 
    
    self.init = function() {
    }

    self.onLoading = function() {
        $.ajax({
            url: "/data/get_evt_endpoints",
        }).done(function(res) {
          var endpoints = res.endpoints;
          if (endpoints.length > 0) {
            self.context.endpoints = endpoints;
            self.goto("main")
          } else {
              self.goto("no-records")
          }
        })
    };

    self.onMain = function() {
        var $table = self.$actionSel().find(".table");

        $table.html("");

        self.context.endpoints.sort(function(a, b) {
          if (a.name > b.name) {
            return 1
           } else {
            return -1
           }
        }).forEach(function(e) {
            // Insert row
            embed("EvtEndpointsListRowView", $table, {
                endpoint: e,
                index: self.context.endpoints.indexOf(e),
                uid: self.context.uid,
            });
        });
    };
}

function EvtEndpointsListRowView() {
    this.base = View; this.base(); var self = this; 
    
    self.onLoading = function() {
      setTimeout(function() {
        self.goto("main")
      }, self.context.index*100)
    };

    self.onMain = function() {
      self.$actionSel().find(".name").html(self.context.endpoint.name)

      // If this has filtering enabled, then show rules
      self.$actionSel().find(".rules-btn-container").addClass('hidden');
      self.$actionSel().find(".unmanaged-container").addClass('hidden');

      if (self.context.endpoint.s2_filter_enabled) {
        self.$actionSel().find(".rules-btn-container").removeClass('hidden');
        self.$actionSel().find(".rules-btn").on("click", function() {
          window.location = "/evt_endpoints/" + self.context.endpoint.name + "/edit_rules";
        })
      } else {
        self.$actionSel().find(".unmanaged-container").removeClass('hidden');
      }
    };
}

function EvtEndpointEditRulesView() {
    this.base = View; this.base(); var self = this; 
    
    self.init = function() {
    }

    self.onLoading = function() {
        $.ajax({
            url: "/data/get_evt_endpoints/" + self.context.endpointName + "/filters",
        }).done(function(res) {
          console.log(res);
            var filters = res.filters;
            if (filters.length > 0) {
              self.context.filters = filters;
              self.goto("main")
            } else {
                self.goto("no-records")
            }
        })
    };

    self.onMain = function() {
        var $tableContent = self.$actionSel().find(".table-content");

        $tableContent.html("");

        self.context.filters.forEach(function(e) {
            // Insert row
            embed("EvtEndpointEditRulesRowView", $tableContent, {
              filter: e,
              index: self.context.filters.indexOf(e),
            });
        });
    };
}

function EvtEndpointEditRulesRowView() {
    this.base = View; this.base(); var self = this; 

    this.init = function() {
        self.$sel(".toggle").on("click", function() {
            self.context.updates = self.context.updates || {};

            self.context.updates['pass'] = $(this).is(':checked');
            self.goto("update", $(this).val())
        })
    }

    self.onUpdate = function() {
      var data = self.context.updates;
      data.filter_id = self.context.filter.id;
        $.ajax({
            method: "put",
            url: "/data/update_evt_filter/",
            contentType: 'application/json',
            data: JSON.stringify(data),
        }).done(function(res) {
          console.log(res);
            self.context.updates = {};
            self.context.filter = res;
            self.goto("main")
        });
    }

    
    self.onLoading = function() {
      setTimeout(function() {
        self.goto("main")
      }, self.context.index*50)
    };

    self.onMain = function() {
      console.log(self.context)
      var $name = self.$actionSel().find(".name")
      var $count = self.$actionSel().find(".count")
      var $toggle = self.$actionSel().find(".toggle")

      var filter = self.context.filter;

      $name.html(filter.event_name);
      $count.html(filter.count);
      $toggle.prop("checked", filter.pass)
    };
}
