$module('compare', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.$equals = function (first, second) {
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      });
      return $promise.build($state);
    };
    $static.$compare = function (first, second) {
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(1, $g.____testlib.basictypes.Integer));
        return;
      });
      return $promise.build($state);
    };
  });

  $static.TEST = function () {
    var first;
    var second;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.compare.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            first = $result;
            $g.compare.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            second = $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $t.box(!$t.unbox($result0), $g.____testlib.basictypes.Boolean);
              $state.current = 4;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) < 0, $g.____testlib.basictypes.Boolean);
              $state.current = 5;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 5:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) > 0, $g.____testlib.basictypes.Boolean);
              $state.current = 6;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 6:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) <= 0, $g.____testlib.basictypes.Boolean);
              $state.current = 7;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 7:
            $result;
            $g.compare.SomeClass.$compare(first, second).then(function ($result0) {
              $result = $t.box($t.unbox($result0) >= 0, $g.____testlib.basictypes.Boolean);
              $state.current = 8;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 8:
            $result;
            $g.compare.SomeClass.$equals(first, second).then(function ($result0) {
              $result = $result0;
              $state.current = 9;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 9:
            $state.resolve($result);
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
