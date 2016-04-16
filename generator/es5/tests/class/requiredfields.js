$module('requiredfields', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField) {
      var instance = new $static();
      var init = [];
      init.push($promise.new(function (resolve) {
        instance.SomeField = SomeField;
        resolve();
      }));
      init.push($promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.AnotherField = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  $static.TEST = function () {
    var sc;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.requiredfields.SomeClass.new($t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $result0;
              $state.current = 1;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            $g.____testlib.basictypes.Integer.$equals(sc.SomeField, $t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
              return $promise.resolve($t.unbox($result1)).then(function ($result0) {
                $result = $t.box($result0 && $t.unbox(sc.AnotherField), $g.____testlib.basictypes.Boolean);
                $state.current = 2;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
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
