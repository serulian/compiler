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
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['AnotherBool'], right.$data['AnotherBool'], $g.____testlib.basictypes.Boolean));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'AnotherBool', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['AnotherBool'], $g.____testlib.basictypes.Boolean, false, 'AnotherBool');
        }
        if (this.$data['AnotherBool'] != null) {
          return $t.box(this.$data['AnotherBool'], $g.____testlib.basictypes.Boolean);
        }
        return this.$data['AnotherBool'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['AnotherBool'] = $t.unbox(val);
          return;
        }
        this.$data['AnotherBool'] = val;
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
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['SomeField'], right.$data['SomeField'], $g.____testlib.basictypes.Integer));
      promises.push($t.equals(left.$data['AnotherField'], right.$data['AnotherField'], $g.____testlib.basictypes.Boolean));
      promises.push($t.equals(left.$data['SomeInstance'], right.$data['SomeInstance'], $g.basic.AnotherStruct));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'SomeField', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['SomeField'], $g.____testlib.basictypes.Integer, false, 'SomeField');
        }
        if (this.$data['SomeField'] != null) {
          return $t.box(this.$data['SomeField'], $g.____testlib.basictypes.Integer);
        }
        return this.$data['SomeField'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['SomeField'] = $t.unbox(val);
          return;
        }
        this.$data['SomeField'] = val;
      },
    });
    Object.defineProperty($instance, 'AnotherField', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['AnotherField'], $g.____testlib.basictypes.Boolean, false, 'AnotherField');
        }
        if (this.$data['AnotherField'] != null) {
          return $t.box(this.$data['AnotherField'], $g.____testlib.basictypes.Boolean);
        }
        return this.$data['AnotherField'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['AnotherField'] = $t.unbox(val);
          return;
        }
        this.$data['AnotherField'] = val;
      },
    });
    Object.defineProperty($instance, 'SomeInstance', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['SomeInstance'], $g.basic.AnotherStruct, false, 'SomeInstance');
        }
        if (this.$data['SomeInstance'] != null) {
          return $t.box(this.$data['SomeInstance'], $g.basic.AnotherStruct);
        }
        return this.$data['SomeInstance'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['SomeInstance'] = $t.unbox(val);
          return;
        }
        this.$data['SomeInstance'] = val;
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
