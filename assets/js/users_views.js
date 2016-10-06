function UserDecodeProfileUrlView() {
    this.base = View; this.base(); var self = this;

    self.init = function() {
        self.$sel(".find-user-btn").click(function() {
          var url = self.$actionSel().find(".url-field").val();
          self.context.url = url;

          self.goto("loading");
        });
    }

    self.onMain = function() {
    }

    self.onLoading = function() {
      // Id is end of url
      var hash = self.context.url.replace(/.*?u=/, "");

      $.ajax({
          method: "post",
          url: "/data/decode_profile_hash",
          contentType: 'application/json',
          data: JSON.stringify({
            hash: hash,
          }),
      }).done(function(res) {
        window.location = '/users/' + res.uid + '/show'
      }).fail(function(err) {
        self.goto("main");
        alert(JSON.stringify(err))
      });
    }
}
