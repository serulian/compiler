$module('withexit', function () {
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
              $g.withexit.someBool = $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
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
            $t.nominalwrap(123, $g.____testlib.basictypes.Integer);
            $state.current = 1;
            continue;

          case 1:
            $g.withexit.SomeReleasable.new().then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $temp0 = $result;
            $state.pushr($temp0, '$temp0');
            $t.nominalwrap(456, $g.____testlib.basictypes.Integer);
            if (false) {
              $state.current = 3;
              continue;
            } else {
              $state.current = 5;
              continue;
            }
            break;

          case 3:
            $state.popr('$temp0').then(function () {
              $state.current = 4;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            continue;

          case 5:
            $t.nominalwrap(012, $g.____testlib.basictypes.Integer);
            $state.popr('$temp0').then(function ($result0) {
              $result = $result0;
              $state.current = 6;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 6:
            $result;
            $state.resolve($g.withexit.someBool);
            return;

          default:
            $state.current = -1;
            return;
        }
      }
    });
    return $promise.build($state);
  };
  this.$init($promise.resolve($t.nominalwrap(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
    $static.someBool = result;
  }));
});
