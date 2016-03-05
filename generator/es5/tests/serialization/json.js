$module('json', function () {
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
          $t.ensurevalue(this.$data.SomeInstance, $g.json.AnotherStruct);
        }
        return $t.box(this.$data.SomeInstance, $g.json.AnotherStruct);
      },
      set: function (val) {
        this.$data.SomeInstance = $t.unbox(val, $g.json.AnotherStruct);
      },
    });
  });

  $static.TEST = function () {
    var correct;
    var jsonString;
    var parsed;
    var s;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.json.AnotherStruct.new($t.nominalwrap(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $temp0 = $result0;
              return $g.json.SomeStruct.new($t.nominalwrap(2, $g.____testlib.basictypes.Integer), $t.nominalwrap(false, $g.____testlib.basictypes.Boolean), ($temp0, $temp0)).then(function ($result1) {
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
            s = $result;
            jsonString = $t.nominalwrap('{"AnotherField":false,"SomeField":2,"SomeInstance":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
            s.Stringify($g.____testlib.basictypes.JSON)().then(function ($result0) {
              return $g.____testlib.basictypes.String.$equals($result0, jsonString).then(function ($result1) {
                $result = $result1;
                $state.current = 2;
                $callback($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            correct = $result;
            $g.json.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            parsed = $result;
            $g.____testlib.basictypes.Integer.$equals(parsed.SomeField, $t.nominalwrap(2, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $result = $t.nominalwrap((($t.nominalunwrap(correct) && $t.nominalunwrap($result0)) && !$t.nominalunwrap(parsed.AnotherField)) && $t.nominalunwrap(parsed.SomeInstance.AnotherBool), $g.____testlib.basictypes.Boolean);
              $state.current = 4;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
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
