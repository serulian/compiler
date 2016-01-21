$module('with', function () {
  var $static = this;
  this.$class('SomeReleasable', false, function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Release = function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.with.someBool = true;
              $state.current = -1;
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    };
  });

  $static.TEST = function () {
    var $temp0;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            123;
            $g.with.SomeReleasable.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $temp0 = $result;
            $state.pushr($temp0, '$temp0');
            456;
            $state.popr('$temp0').then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $result;
            789;
            $state.resolve($g.with.someBool);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
  this.$init($promise.resolve(false).then(function (result) {
    $static.someBool = result;
  }));
});
