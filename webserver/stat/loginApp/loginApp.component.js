(function(loginApp) {
  loginApp.AppComponent =
    ng.core.Component({
      selector: 'login-app',
      templateUrl: 'stat/loginApp/template.html'
    })
    .Class({
      constructor: [ ng.http.Http, function(http) {
        this.http = http;
      }],

      login: function() {
        this.http
          .post("login", { username : this.username, password : this.password })
          .map(function(result) {
            return result.json()
          })
          .subscribe(function(success) {
            if (success) {
              window.location.href = "/"
            } else {
              this.wrongCredentials = true;
            }
          }.bind(this));
      }
    });
})(window.loginApp || (window.loginApp = {}));