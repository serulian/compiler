$module('basic', function () {
  var $static = this;
  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.AnotherBool = AnotherBool;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['AnotherBool'] = this.AnotherBool;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    Object.defineProperty($instance, 'AnotherBool', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data.AnotherBool, $g.____testlib.basictypes.Boolean);
        }
        return $t.box(this.$data.AnotherBool, $g.____testlib.basictypes.Boolean);
      },
      set: function (val) {
        this.$data.AnotherBool = $t.unbox(val, $g.____testlib.basictypes.Boolean);
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
      instance.$lazycheck = false;
      instance.SomeField = SomeField;
      instance.AnotherField = AnotherField;
      instance.SomeInstance = SomeInstance;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['SomeField'] = this.SomeField;
      mappedData['AnotherField'] = this.AnotherField;
      mappedData['SomeInstance'] = this.SomeInstance;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data.SomeField, $g.____testlib.basictypes.Integer);
        }
        return $t.box(this.$data.SomeField, $g.____testlib.basictypes.Integer);
      },
      set: function (val) {
        this.$data.SomeField = $t.unbox(val, $g.____testlib.basictypes.Integer);
      },
    });
    Object.defineProperty($instance, 'AnotherField', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data.AnotherField, $g.____testlib.basictypes.Boolean);
        }
        return $t.box(this.$data.AnotherField, $g.____testlib.basictypes.Boolean);
      },
      set: function (val) {
        this.$data.AnotherField = $t.unbox(val, $g.____testlib.basictypes.Boolean);
      },
    });
    Object.defineProperty($instance, 'SomeInstance', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data.SomeInstance, $g.basic.AnotherStruct);
        }
        return $t.box(this.$data.SomeInstance, $g.basic.AnotherStruct);
      },
      set: function (val) {
        this.$data.SomeInstance = $t.unbox(val, $g.basic.AnotherStruct);
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
