$module('slice', function () {
  var $static = this;
  this.$struct('AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherInt) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        AnotherInt: AnotherInt,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'AnotherInt', {
        get: function () {
          if (this.$lazychecked['AnotherInt']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['AnotherInt'], $g.____testlib.basictypes.Integer, false, 'AnotherInt');
            this.$lazychecked['AnotherInt'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['AnotherInt'], $g.____testlib.basictypes.Integer);
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['AnotherInt'] = this.AnotherInt;
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
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['AnotherInt'], right[BOXED_DATA_PROPERTY]['AnotherInt'], $g.____testlib.basictypes.Integer));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'AnotherInt', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['AnotherInt'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['AnotherInt'] = value;
      },
    });
  });

  this.$struct('SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (Values) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        Values: Values,
      };
      return $promise.resolve(instance);
    };
    $static.$box = function (data) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = data;
      instance.$lazychecked = {
      };
      Object.defineProperty(instance, 'Values', {
        get: function () {
          if (this.$lazychecked['Values']) {
            $t.ensurevalue(this[BOXED_DATA_PROPERTY]['Values'], $g.____testlib.basictypes.Slice($g.slice.AnotherStruct), false, 'Values');
            this.$lazychecked['Values'] = true;
          }
          return $t.box(this[BOXED_DATA_PROPERTY]['Values'], $g.____testlib.basictypes.Slice($g.slice.AnotherStruct));
        },
      });
      instance.Mapping = function () {
        var mapped = {
        };
        mapped['Values'] = this.Values;
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
      promises.push($t.equals(left[BOXED_DATA_PROPERTY]['Values'], right[BOXED_DATA_PROPERTY]['Values'], $g.____testlib.basictypes.Slice($g.slice.AnotherStruct)));
      return Promise.all(promises).then(function (values) {
        for (var i = 0; i < values.length; i++) {
          if (!$t.unbox(values[i])) {
            return values[i];
          }
        }
        return $t.box(true, $g.____testlib.basictypes.Boolean);
      });
    };
    Object.defineProperty($instance, 'Values', {
      get: function () {
        return this[BOXED_DATA_PROPERTY]['Values'];
      },
      set: function (value) {
        this[BOXED_DATA_PROPERTY]['Values'] = value;
      },
    });
  });

  $static.TEST = function () {
    var correct;
    var jsonString;
    var parsed;
    var s;
    var values;
    var $state = $t.sm(function ($continue) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.slice.AnotherStruct.new($t.box(1, $g.____testlib.basictypes.Integer)).then(function ($result0) {
              $temp0 = $result0;
              return $g.slice.AnotherStruct.new($t.box(2, $g.____testlib.basictypes.Integer)).then(function ($result1) {
                $temp1 = $result1;
                return $g.slice.AnotherStruct.new($t.box(3, $g.____testlib.basictypes.Integer)).then(function ($result2) {
                  $temp2 = $result2;
                  return $g.____testlib.basictypes.List($g.slice.AnotherStruct).forArray([($temp0, $temp0), ($temp1, $temp1), ($temp2, $temp2)]).then(function ($result3) {
                    $result = $result3;
                    $state.current = 1;
                    $continue($state);
                  });
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 1:
            values = $result;
            values.$slice($t.box(0, $g.____testlib.basictypes.Integer), null).then(function ($result0) {
              return $g.slice.SomeStruct.new($result0).then(function ($result1) {
                $temp3 = $result1;
                $result = ($temp3, $temp3);
                $state.current = 2;
                $continue($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 2:
            s = $result;
            jsonString = $t.box('{"Values":[{"AnotherInt":1},{"AnotherInt":2},{"AnotherInt":3}]}', $g.____testlib.basictypes.String);
            s.Stringify($g.____testlib.basictypes.JSON)().then(function ($result0) {
              return $g.____testlib.basictypes.String.$equals($result0, jsonString).then(function ($result1) {
                $result = $result1;
                $state.current = 3;
                $continue($state);
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 3:
            correct = $result;
            $g.slice.SomeStruct.Parse($g.____testlib.basictypes.JSON)(jsonString).then(function ($result0) {
              $result = $result0;
              $state.current = 4;
              $continue($state);
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 4:
            parsed = $result;
            $promise.resolve($t.unbox(correct)).then(function ($result1) {
              return ($promise.shortcircuit($result1, false) || s.Values.Length()).then(function ($result2) {
                return ($promise.shortcircuit($result1, false) || $g.____testlib.basictypes.Integer.$equals($result2, $t.box(3, $g.____testlib.basictypes.Integer))).then(function ($result3) {
                  return $promise.resolve($result1 && $t.unbox($result3)).then(function ($result0) {
                    return ($promise.shortcircuit($result0, false) || s.Values.$index($t.box(0, $g.____testlib.basictypes.Integer))).then(function ($result4) {
                      return ($promise.shortcircuit($result0, false) || $g.____testlib.basictypes.Integer.$equals($result4.AnotherInt, $t.box(1, $g.____testlib.basictypes.Integer))).then(function ($result5) {
                        $result = $t.box($result0 && $t.unbox($result5), $g.____testlib.basictypes.Boolean);
                        $state.current = 5;
                        $continue($state);
                      });
                    });
                  });
                });
              });
            }).catch(function (err) {
              $state.reject(err);
            });
            return;

          case 5:
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
