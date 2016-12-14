(function(mainApp) {
  mainApp.UsersComponent =
    ng.core.Component({
      selector: 'users-cpt',
      templateUrl: 'stat/mainApp/users/template.html'
    })
    .Class({
      constructor: [ ng.http.Http, function(http) {
        this.http = http;
        this.http
          .get("api/users")
          .map(function(result) {
            return result.json()
          })
          .subscribe(function(result) {
              result.sort(function(a, b) {
                  return a.usertype.localeCompare(b.usertype);
              });

            this.users = result;
          }.bind(this));
      }],

      trackFunc: function(index, user) {
          return user.username;
      }
    });
})(window.mainApp || (window.mainApp = {}));