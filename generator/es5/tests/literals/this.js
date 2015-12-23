$module('this', function () {
  var $static = this;
  this.cls('SomeClass', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function ($callback) {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.DoSomething = function () {
      var $this = this;
      var $state = {
        current: 0,
        returnValue: null,
      };
      $state.next = function ($callback) {
        try {
          while (true) {
            switch ($state.current) {
              case 0:
                $this;
                $state.current = -1;
                return;

              default:
                $state.current = -1;
                return;
            }
          }
        } catch (e) {
          $state.error = e;
          $state.current = -1;
          $callback($state);
        }
      };
      return $promise.build($state);
    };
  });

});
