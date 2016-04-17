$module('nominaljson', function () {
  var $static = this;
  this.$type('SomeNominal', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    $instance.GetValue = function () {
      var $this = this;
      var $state = $t.sm(function ($continue) {
        $state.resolve($t.box($this, $g.nominaljson.AnotherStruct).AnotherBool);
        return;
      });
      return $promise.build($state);
    };
  });

  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        AnotherBool: AnotherBool,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'AnotherBool', {
        get: function () {
          if (this.$lazychecked['AnotherBool']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['AnotherBool'], $g.____testlib.basictypes.Boolean, false, 'AnotherBool');
            this.$lazychecked['AnotherBool'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['AnotherBool'], $g.____testlib.basictypes.Boolean);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['AnotherBool'] = this.AnotherBool;
        return $promise.resolve($t.box(mapped, $g.____testlib.basictypes.Mapping($t.any)));
      };
      return instance;
    };
    $instance.Mapping = function () {
      return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['AnotherBool'], right[BOXED_DATA_PROPERTY]['AnotherBool'], $g.____testlib.basictypes.Boolean));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'AnotherBool', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['AnotherBool'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['AnotherBool'] = value;
      },
    });
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Nested) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Nested: Nested,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'Nested', {
        get: function () {
          if (this.$lazychecked['Nested']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['Nested'], $g.nominaljson.SomeNominal, false, 'Nested');
            this.$lazychecked['Nested'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['Nested'], $g.nominaljson.SomeNominal);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['Nested'] = this.Nested;
        return $promise.resolve($t.box(mapped, $g.____testlib.basictypes.Mapping($t.any)));
      };
      return instance;
    };
    $instance.Mapping = function () {
      return $promise.resolve($t.box(this[BOXED_DATA_PROPERTY], $g.____testlib.basictypes.Mapping($t.any)));
    };
    $static.$equals = function (left, right) {
      if (left === right) {
        return $promise.resolve($t.box(true, $g.____testlib.basictypes.Boolean));
      }
      var promises = [];
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['Nested'], right[BOXED_DATA_PROPERTY]['Nested'], $g.nominaljson.SomeNominal));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'Nested', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['Nested'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['Nested'] = value;
      },
    });
  });

  $static.TEST = function () {
    var correct;
    var jsonString;
    var parsed;
    var s;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.nominaljson.AnotherStruct.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $temp0 = $result0;
              return $g.nominaljson.SomeStruct.new($t.box(($temp0, $temp0), $g.nominaljson.SomeNominal)).then(function ($result1) {
                $temp1 = $result1;
                $result = ($temp1, $temp1);
                $state.current = 1;
                $continue($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            s = $result;
            jsonString = $t.box('{"Nested":{"AnotherBool":true}}', $g.____testlib.basictypes.String);
            s.Stringify($g.____testlib.basictypes.JSON)().then(function ($result0) {
              return $g.____testlib.basictypes.String.$equals($result0, jsonString).then(function ($result1) {
                $result = $result1;
                $state.current = 2;
                $continue($state);
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
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            parsed = $result;
            $promise.resolve($t.unbox(correct)).then(function ($result0) {
              return ($promise.shortcircuit($result0, false) || parsed.Nested.GetValue()).then(function ($result1) {
                $result = $t.box($result0 && $t.unbox($result1), $g.____testlib.basictypes.Boolean);
                $state.current = 4;
                $continue($state);
              });
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
