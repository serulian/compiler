$module('basic', function () {
  var $static = this;
  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance.$data = {
      };
      instance.AnotherBool = AnotherBool;
      return $promise.resolve(instance);
    };
    $static.$apply = function (toplevel) {
      var data = toplevel.$data;
      var instance = new $static();
      instance.$data = {
      };
      instance.AnotherBool = $g.____testlib.basictypes.Boolean.$apply(data['AnotherBool']);
      return instance;
    };
    Object.defineProperty($instance, 'AnotherBool', {
      get: function () {
        return this.$data.AnotherBool;
      },
      set: function (val) {
        this.$data.AnotherBool = val;
      },
    });
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField, SomeInstance) {
      var instance = new $static();
      instance.$data = {
      };
      instance.SomeField = SomeField;
      instance.AnotherField = AnotherField;
      instance.SomeInstance = SomeInstance;
      return $promise.resolve(instance);
    };
    $static.$apply = function (toplevel) {
      var data = toplevel.$data;
      var instance = new $static();
      instance.$data = {
      };
      instance.SomeField = $g.____testlib.basictypes.Integer.$apply(data['SomeField']);
      instance.AnotherField = $g.____testlib.basictypes.Boolean.$apply(data['AnotherField']);
      instance.SomeInstance = $g.basic.AnotherStruct.$apply(data['SomeInstance']);
      return instance;
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        return this.$data.SomeField;
      },
      set: function (val) {
        this.$data.SomeField = val;
      },
    });
    Object.defineProperty($instance, 'AnotherField', {
      get: function () {
        return this.$data.AnotherField;
      },
      set: function (val) {
        this.$data.AnotherField = val;
      },
    });
    Object.defineProperty($instance, 'SomeInstance', {
      get: function () {
        return this.$data.SomeInstance;
      },
      set: function (val) {
        this.$data.SomeInstance = val;
      },
    });
  });

  $static.TEST = function () {
    var ss;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.basic.AnotherStruct.new($t.nominalwrap(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $temp0 = $result0;
              return $g.basic.SomeStruct.new($t.nominalwrap(42, $g.____testlib.basictypes.Integer), $t.nominalwrap(true, $g.____testlib.basictypes.Boolean), ($temp0, $temp0)).then(function ($result1) {
                $temp1 = $result1;
                $result = ($temp1, $temp1);
                $state.current = 1;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            ss = $result;
            $g.____testlib.basictypes.Integer.$equals(ss.SomeField, $t.nominalwrap(42, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $t.nominalwrap(($t.nominalunwrap($result0) && $t.nominalunwrap(ss.AnotherField)) && $t.nominalunwrap(ss.SomeInstance.AnotherBool), $g.____testlib.basictypes.Boolean);
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
