$(function() {
  RegisterTemplates();
  LoadInlineEmbeds();
});

templateEngine = {
  templates: {},
  views: {},
}

function RegisterTemplates() {
    //////////////////////////////////////////////////
    // Register templates
    //////////////////////////////////////////////////
    $(".view-templates .view-template").each(function() {
        var templateInfo = {};

        // Get name of template
        var name = $(this).attr('data-name');

        // Bind HTML to template
        templateInfo.$sel = $(this);

        // Lookup constructor class for template
        templateInfo.constructor = window[name];

        console.log("Registered: " + name)
        templateEngine.templates[name] = templateInfo;
    });
    //////////////////////////////////////////////////
}

function LoadInlineEmbeds() {
    //////////////////////////////////////////////////
    // Our root views will be embedded directly by searching the HTML
    // for divs with data-embeds containing the view name.
    //////////////////////////////////////////////////
    $(".embed").each(function() {
        var templateName = $(this).attr('data-name');
        var _this = this;
        console.log("embedding " + templateName);

        // Schedule off the main run-queue
        setTimeout(function() {
          embed(templateName, $(_this));
        }, 0)
    });
}

// Embed a new view into a JQuery selector location
function embed(templateName, $target, context) {
    context = context || {};

    // Get template information
    var tmpl = templateEngine.templates[templateName];

    // Bind UUID
    var uuid = UUID();

    // Load template at location with UUID
    if (tmpl == null) {
        console.error("Template named: '" + templateName + "' did not exist");
        return
    }

    $target.append("<div class='view templated-view " + templateName + "' data-uuid='" + uuid + "'>" + tmpl.$sel.html() + "</div>");

    // If you have a div like <div class='embed' data-name='xxx' data-context-uid='44'></div> the
    // context.uid of this view will be 44.
    var data = $target.data()
    for (var key in data) {
      var re = /^context/
      if (re.test(key)) {
        var newKey = key.replace(re, "").toCamelCase().lowercaseFirstLetter()
        context[newKey] = data[key];
      }
    }

    // Get new view's $sel
    var $sel = $("[data-uuid='" + uuid + "']");

    // Get any context setters

    rebindToSelector($sel);

    // Retrieve the 'first action' which we'll consider the default
    var firstActionName = $sel.find(".action:nth(0)").attr('data-name');

    // Save information for this live view
    if (tmpl.constructor == null || tmpl.constructor == undefined) {
      console.error("Fatal: Tried to get view controller for: " + templateName + " but this didn't exist");
      return
    }
    
    var view = new tmpl.constructor()
    view.__initialize__(templateName, $sel, context);
    //templateEngine.views[uuid] = view;
    $sel.data("__view__", view)

    view.init()
    view.goto(firstActionName);
}

// View class
function View() {
    this.__initialize__ = function(templateName, $sel, context) {
        this.$_sel = $sel;
        this.templateName = templateName;
        this.context = context;
        this.currentActionName = "";
        this.notificationHandlers = {};
    }

    this.onNotification = function(name, handler) {
      this.notificationHandlers[name] = handler;
    }

    this.$sel = function(str) {
        return $(str, this.$_sel)
    }

    this.$actionSel = function() {
        $res = $(".action[data-name='" + this.currentActionName + "']", this.$_sel).not($(".action *", this.$_sel));
        return $res
    }

    this.init = function() {}
    
    // Notify all views
    this.postNotification = function(name) {
      $(document).find(".view[data-uuid]").each(function() {
        var _this = this;
        setTimeout(function() {
          var v = $(_this).data().__view__;
          var nh = v.notificationHandlers[name];
          if (nh !== undefined && nh !== null) {
            nh(v);
          }
        }, 0);
      });
    }

    // Switch actions
    this.goto = function(actionName) {
        this.currentActionName = actionName;
        this.$sel(".action").removeClass('hidden');
        //this.$sel(".action").animateCss("zoomIn")
        this.$sel(".action").not(".action[data-name='" + actionName + "']").addClass('hidden');

        if (this.currentActionName == null || this.currentActionName == undefined) {
          console.error("Your HTML template named " + this.templateName + " did not contain a divider for the action named: " + this.currentActionName);
        }
        var sel = "on" + this.currentActionName.toClassCase();
        var action = this[sel];
        if (action != null) {
            var args = Array.apply(null, arguments).slice(1);
            action(args);
        } else {}
    }
}

// Auto-rebind selectors for templates
function rebindToSelector($sel) {
    $sel.find(".hero-minimize").on("click", function() {
        $(this).closest('.hero').toggleClass('hero-minimized')
    });

    $sel.find(".toggle + label").on("click", function() {
        $(this).prev().click()
    });

    jQuery.timeago.settings.allowFuture = true;
    $sel.find("time.timeago").timeago();
};


