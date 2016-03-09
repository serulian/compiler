$module('custom', function () {
  var $static = this;
  this.$class('CustomJSON', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $static.Get = function () {
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.custom.CustomJSON.new().then(function ($result0) {
                $result = $result0;
                $state.current = 1;
                $callback($state);
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
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
    $instance.Stringify = function (value) {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.____testlib.basictypes.JSON.Get().then(function ($result0) {
                return $result0.Stringify(value).then(function ($result1) {
                  $result = $result1;
                  $state.current = 1;
                  $callback($state);
                });
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
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
    $instance.Parse = function (value) {
      var $this = this;
      var $state = $t.sm(function ($callback) {
        while (true) {
          switch ($state.current) {
            case 0:
              $g.____testlib.basictypes.JSON.Get().then(function ($result0) {
                return $result0.Parse(value).then(function ($result1) {
                  $result = $result1;
                  $state.current = 1;
                  $callback($state);
                });
              }).catch(function (err) {
                $state.reject(err);
              });
              return;

            case 1:
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
          $t.ensurevalue(this.$data['SomeInstance'], $g.custom.AnotherStruct, false, 'SomeInstance');
        }
        if (this.$data['SomeInstance'] != null) {
          return $t.box(this.$data['SomeInstance'], $g.custom.AnotherStruct);
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
    var correct;
    var jsonString;
    var parsed;
    var s;
    var $state = $t.sm(function ($callback) {
      while (true) {
        switch ($state.current) {
          case 0:
            $g.custom.AnotherStruct.new($t.nominalwrap(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $temp0 = $result0;
              return $g.custom.SomeStruct.new($t.nominalwrap(2, $g.____testlib.basictypes.Integer), $t.nominalwrap(false, $g.____testlib.basictypes.Boolean), ($temp0, $temp0)).then(function ($result1) {
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
            s.Stringify($g.custom.CustomJSON)().then(function ($result0) {
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
            $g.custom.SomeStruct.Parse($g.custom.CustomJSON)(jsonString).then(function ($result0) {
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
