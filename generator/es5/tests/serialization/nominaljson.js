$module('nominaljson', function () {
  var $static = this;
  this.$type('SomeNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.new = function ($wrapped) {
      var instance = new this();
      instance.$wrapped = $wrapped;
      return instance;
    };
    this.$box = function (data) {
      var instance = new this();
      instance.$wrapped = $t.box(data, $g.nominaljson.AnotherStruct);
      return instance;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $state.resolve($t.nominalunwrap($this).AnotherBool);
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
    $static.new = function (Nested) {
      var instance = new $static();
      instance.$data = {
      };
      instance.$lazycheck = false;
      instance.Nested = Nested;
      return $promise.resolve(instance);
    };
    $instance.Mapping = function () {
      var mappedData = {
      };
      mappedData['Nested'] = this.Nested;
      return $promise.resolve($t.nominalwrap(mappedData, $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.nominalwrap(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left.$data['Nested'], right.$data['Nested'], $g.nominaljson.SomeNominal));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return $t.nominalwrap(false, $g.____testlib.basictypes.Boolean);
          }
        }
        return $t.nominalwrap(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'Nested', {
      get: function () {
        if (this.$lazycheck) {
          $t.ensurevalue(this.$data['Nested'], $g.nominaljson.SomeNominal, false, 'Nested');
        }
        if (this.$data['Nested'] != null) {
          return $t.box(this.$data['Nested'], $g.nominaljson.SomeNominal);
        }
        return this.$data['Nested'];
      },
      set: function (val) {
        if (val != null) {
          this.$data['Nested'] = $t.unbox(val);
          return;
        }
        this.$data['Nested'] = val;
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
            $g.nominaljson.AnotherStruct.new($t.nominalwrap(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $temp0 = $result0;
              return $g.nominaljson.SomeStruct.new($t.nominalwrap(($temp0, $temp0), $g.nominaljson.SomeNominal)).then(function ($result1) {
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
            jsonString = $t.nominalwrap('{"Nested":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
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
            $g.nominaljson.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
              $result = $result0;
              $state.current = 3;
              $callback($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            parsed = $result;
            parsed.Nested.GetValue().then(function ($result0) {
              $result = $t.nominalwrap($t.nominalunwrap(correct) && $t.nominalunwrap($result0), $g.____testlib.basictypes.Boolean);
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
