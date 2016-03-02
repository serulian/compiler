$module('basic', function () {
  var $static = this;
  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance.data = {
      };
      instance.AnotherBool = AnotherBool;
      return $promise.resolve(instance);
    };
    Object.defineProperty($instance, 'AnotherBool', {
      get: function () {
        return $t.nominalwrap(this.data.AnotherBool, $g.____testlib.basictypes.Boolean);
      },
      set: function (val) {
        this.data.AnotherBool = $t.nominalunwrap(val);
      },
    });
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField, SomeInstance) {
      var instance = new $static();
      instance.data = {
      };
      instance.SomeField = SomeField;
      instance.AnotherField = AnotherField;
      instance.SomeInstance = SomeInstance;
      return $promise.resolve(instance);
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        return $t.nominalwrap(this.data.SomeField, $g.____testlib.basictypes.Integer);
      },
      set: function (val) {
        this.data.SomeField = $t.nominalunwrap(val);
      },
    });
    Object.defineProperty($instance, 'AnotherField', {
      get: function () {
        return $t.nominalwrap(this.data.AnotherField, $g.____testlib.basictypes.Boolean);
      },
      set: function (val) {
        this.data.AnotherField = $t.nominalunwrap(val);
      },
    });
    Object.defineProperty($instance, 'SomeInstance', {
      get: function () {
        return this.data.SomeInstance;
      },
      set: function (val) {
        this.data.SomeInstance = val;
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
