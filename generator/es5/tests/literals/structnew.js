$module('structnew', function () {
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
      init.push($promise.resolve($t.nominalwrap(false, $g.____testlib.basictypes.Boolean)).then(function (result) {
        instance.anotherField = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.AnotherField = $t.property(function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($this.anotherField);
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    }, function (val) {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $this.anotherField = val;
              $state.current = -1;
              return;

            default:
              $state.current = -1;
              return;
          }
        }
      });
      return $promise.build($state);
    });
  });

  $static.TEST = function () {
    var sc;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.structnew.SomeClass.new($t.nominalwrap(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $temp0 = $result0;
              return $temp0.AnotherField($t.nominalwrap(true, $g.____testlib.basictypes.Boolean)).then(function ($result1) {
                $result = ($temp0, $result1, $temp0);
                $state.current = 1;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            sc = $result;
            $g.____testlib.basictypes.Integer.$equals(sc.SomeField, $t.nominalwrap(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $t.nominalwrap($t.nominalunwrap($result0) && $t.nominalunwrap(sc.anotherField), $g.____testlib.basictypes.Boolean);
              $state.current = 2;
              $callback($state);
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
