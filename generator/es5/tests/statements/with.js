$module('with', function () {
  var $static = this;
  this.$class('SomeReleasable', false, '', function () {
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
      var $state = $t.sm(function ($continue) {
        $g.with.someBool = $t.box(true, $g.____testlib.basictypes.Boolean);
        $state.resolve();
      });
      return $promise.build($state);
    };
  });

  $static.TEST = function () {
    var $temp0;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $t.box(123, $g.____testlib.basictypes.Integer);
            $g.with.SomeReleasable.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            $temp0 = $result;
            $state.pushr($temp0, '$temp0');
            $t.box(456, $g.____testlib.basictypes.Integer);
            $state.popr('$temp0').then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $result;
            $t.box(789, $g.____testlib.basictypes.Integer);
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
  this.$init(function () {
    return $promise.resolve($t.box(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
      $static.someBool = result;
    });
  });
});
