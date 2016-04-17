$module('slice', function () {
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
    $instance.$slice = function (start, end) {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      });
      return $promise.build($state);
    };
  });

  $static.TEST = function () {
    var c;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.slice.SomeClass.new().then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            c = $result;
            c.$slice($t.box(1, $g.____testlib.basictypes.Integer), $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $state.current = 2;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            $result;
            c.$slice(null, $t.box(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            $result;
            c.$slice($t.box(1, $g.____testlib.basictypes.Integer), null).then(function ($result0) {
              $result = $result0;
              $state.current = 4;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
            $result;
            c.$slice($t.box(1, $g.____testlib.basictypes.Integer), $t.box(7, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $state.current = 5;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 5:
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
