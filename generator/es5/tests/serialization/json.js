$module('json', function () {
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
      instance.SomeInstance = $g.json.AnotherStruct.$apply(data['SomeInstance']);
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
